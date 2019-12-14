# General imports
import os
import sys
from multiprocessing.pool import ThreadPool

import tornado.gen
import tornado.httpserver
import tornado.ioloop
import tornado.web
from raven.contrib.tornado import AsyncSentryClient
import redis

import json
import shutil
from distutils.version import LooseVersion

from constants import rankedStatuses

from common.constants import bcolors, mods
from common.db import dbConnector
from common.ddog import datadogClient
from common.log import logUtils as log
from common.redis import pubSub
from common.web import schiavo
from handlers import apiCacheBeatmapHandler, rateHandler, changelogHandler
from handlers import apiPPHandler
from handlers import apiStatusHandler
from handlers import banchoConnectHandler
from handlers import checkUpdatesHandler
from handlers import defaultHandler
from handlers import downloadMapHandler
from handlers import emptyHandler
from handlers import getFullReplayHandler
from handlers import getFullReplayHandlerRelax
from handlers import getReplayHandler
from handlers import getScoresHandler
from handlers import getScreenshotHandler
from handlers import loadTestHandler
from handlers import mapsHandler
from handlers import inGameRegistrationHandler
from handlers import getFullErrorHandler
from handlers import osuErrorHandler
from handlers import osuSearchHandler
from handlers import osuSearchSetHandler
from handlers import redirectHandler
from handlers import submitModularHandler
from handlers import uploadScreenshotHandler
from handlers import commentHandler
from helpers import config
from helpers import consoleHelper
from common import generalUtils
from common import agpl
from objects import glob
from pubSubHandlers import beatmapUpdateHandler
import secret.achievements.utils
from subprocess import call


def make_app():
	return tornado.web.Application([
		(r"/users", inGameRegistrationHandler.handler),
		(r"/web/bancho_connect.php", banchoConnectHandler.handler),
		(r"/web/osu-osz2-getscores.php", getScoresHandler.handler),
		(r"/web/osu-submit-modular.php", submitModularHandler.handler),
		(r"/web/osu-submit-modular-selector.php", submitModularHandler.handler),
		(r"/web/osu-getreplay.php", getReplayHandler.handler),
		(r"/web/osu-screenshot.php", uploadScreenshotHandler.handler),
		(r"/web/osu-search.php", osuSearchHandler.handler),
		(r"/web/osu-search-set.php", osuSearchSetHandler.handler),
		(r"/web/check-updates.php", checkUpdatesHandler.handler),
		(r"/web/osu-error.php", osuErrorHandler.handler),
		(r"/web/osu-comment.php", commentHandler.handler),
		(r"/p/changelog", changelogHandler.handler),
		(r"/web/changelog.php", changelogHandler.handler),
		(r"/home/changelog", changelogHandler.handler),
		(r"/web/osu-rate.php", rateHandler.handler),
		(r"/ss/(.*)", getScreenshotHandler.handler),
		(r"/web/maps/(.*)", mapsHandler.handler),
		(r"/d/(.*)", downloadMapHandler.handler),
		(r"/s/(.*)", downloadMapHandler.handler),
		(r"/web/replays/(.*)", getFullReplayHandler.handler),
		(r"/web/replays_relax/(.*)", getFullReplayHandlerRelax.handler),
		(r"/web/errorlogs/(.*)", getFullErrorHandler.handler),

		(r"/p/verify", redirectHandler.handler, dict(destination="https://ripple.moe/")),
		(r"/u/(.*)", redirectHandler.handler, dict(destination="https://ripple.moe/u/{}")),

		(r"/api/v1/status", apiStatusHandler.handler),
		(r"/api/v1/pp", apiPPHandler.handler),
		(r"/api/v1/cacheBeatmap", apiCacheBeatmapHandler.handler),

		(r"/letsapi/v1/status", apiStatusHandler.handler),
		(r"/letsapi/v1/pp", apiPPHandler.handler),
		(r"/letsapi/v1/cacheBeatmap", apiCacheBeatmapHandler.handler),

		# Not done yet
		(r"/web/osu-addfavourite.php", osuErrorHandler.handler), # I use the error handler as I want an empty response for the time being
		(r"/web/lastfm.php", emptyHandler.handler),
		(r"/web/osu-checktweets.php", emptyHandler.handler),

		(r"/loadTest", loadTestHandler.handler),
	], default_handler_class=defaultHandler.handler)


if __name__ == "__main__":
    	# AGPL license agreement
	try:
			call('cls', shell=True)
			agpl.check_license("ripple", "LETS")
	except agpl.LicenseError as e:
			print(str(e))
			sys.exit(1)

	try:
		consoleHelper.printServerStartHeader(True)

		# Read config
		consoleHelper.printNoNl("> Reading config file [/]             ")
		glob.conf = config.config("config.ini")

		if glob.conf.default:
			# We have generated a default config.ini, quit server
			consoleHelper.printWarning()
			consoleHelper.printColored("[!] config.ini not found. A default one has been generated.", bcolors.YELLOW)
			consoleHelper.printColored("[!] Please edit your config.ini and run the server again.", bcolors.YELLOW)
			sys.exit()

		# If we haven't generated a default config.ini, check if it's valid
		if not glob.conf.checkConfig():
			consoleHelper.printError()
			consoleHelper.printColored("[!] Invalid config.ini. Please configure it properly", bcolors.RED)
			consoleHelper.printColored("[!] Delete your config.ini to generate a default one", bcolors.RED)
			sys.exit()
		else:
			consoleHelper.printDone()

		# Read additional config file
		consoleHelper.printNoNl("> Loading additional config file [/]  ")
		try:
			if not os.path.isfile(glob.conf.config["custom"]["config"]):
				consoleHelper.printWarning()
				consoleHelper.printColored("[!] Missing config file at {}; A default one has been generated at this location.".format(glob.conf.config["custom"]["config"]), bcolors.YELLOW)
				shutil.copy("common/default_config.json", glob.conf.config["custom"]["config"])

			with open(glob.conf.config["custom"]["config"], "r") as f:
				glob.conf.extra = json.load(f)

			consoleHelper.printDone()
		except:
			consoleHelper.printWarning()
			consoleHelper.printColored("[!] Unable to load custom config at {}".format(glob.conf.config["custom"]["config"]), bcolors.RED)
			consoleHelper.printColored("[!] Make sure you have the latest osu!thailand common submodule!", bcolors.RED)
			sys.exit()

		# Check if running common module is usable
		if glob.COMMON_VERSION == "Unknown":
			consoleHelper.printWarning()
			consoleHelper.printColored("[!] You do not seem to be using osu!thailand's common submodule... nothing will work...", bcolors.RED)
			consoleHelper.printColored("[!] You can download or fork the submodule from {}https://github.com/osuthailand/common".format(bcolors.UNDERLINE), bcolors.RED)
			sys.exit()
		elif LooseVersion(glob.COMMON_VERSION_REQ) > LooseVersion(glob.COMMON_VERSION):
			consoleHelper.printColored("[!] Your common submodule version is below the required version number for this version of lets.", bcolors.RED)
			consoleHelper.printColored("[!] You are highly adviced to update your common submodule as stability may vary with outdated modules.", bcolors.RED)

		# Create data/oppai maps folder if needed
		consoleHelper.printNoNl("> Checking folders [/]                ")
		paths = [
			".data",
			"{}_relax".format(glob.conf.config["server"]["replayspath"]),
			glob.conf.config["server"]["screenshotspath"],
			".data/screenshots",
			".data/oppai",
			".data/catch_the_pp",
			glob.conf.config["server"]["beatmapspath"]
		]
		for i in paths:
			if not os.path.exists(i):
				os.makedirs(i, 0o770)
		consoleHelper.printDone()

		# Connect to db
		try:
			consoleHelper.printNoNl("> Connecting to MySQL database [/]   ")
			glob.db = dbConnector.db(glob.conf.config["db"]["host"], glob.conf.config["db"]["username"], glob.conf.config["db"]["password"], glob.conf.config["db"]["database"], int(
				glob.conf.config["db"]["workers"]))
			consoleHelper.printNoNl(" ")
			consoleHelper.printDone()
		except:
			# Exception while connecting to db
			consoleHelper.printError()
			consoleHelper.printColored("[!] Error while connection to database. Please check your config.ini and run the server again", bcolors.RED)
			raise

		# Connect to redis
		try:
			consoleHelper.printNoNl("> Connecting to Redis [/]            ")
			glob.redis = redis.Redis(glob.conf.config["redis"]["host"], glob.conf.config["redis"]["port"], glob.conf.config["redis"]["database"], glob.conf.config["redis"]["password"])
			glob.redis.ping()
			consoleHelper.printNoNl(" ")
			consoleHelper.printDone()
		except:
			# Exception while connecting to db
			consoleHelper.printError()
			consoleHelper.printColored("[!] Error while connection to redis. Please check your config.ini and run the server again", bcolors.RED)
			raise

		# Empty redis cache
		try:
			glob.redis.eval("return redis.call('del', unpack(redis.call('keys', ARGV[1])))", 0, "lets:*")
		except redis.exceptions.ResponseError:
			# Script returns error if there are no keys starting with peppy:*
			pass

		# Save lets version in redis
		glob.redis.set("lets:version", glob.VERSION)

		# Create threads pool
		try:
			consoleHelper.printNoNl("> Creating threads pool [/]           ")
			glob.pool = ThreadPool(int(glob.conf.config["server"]["threads"]))
			consoleHelper.printDone()
		except:
			consoleHelper.printError()
			consoleHelper.printColored("[!] Error while creating threads pool. Please check your config.ini and run the server again", bcolors.RED)

		# Check osuapi
		if not generalUtils.stringToBool(glob.conf.config["osuapi"]["enable"]):
			consoleHelper.printColored("[!] osu!api features are disabled. If you don't have a valid beatmaps table, all beatmaps will show as unranked", bcolors.YELLOW)
			if int(glob.conf.config["server"]["beatmapcacheexpire"]) > 0:
				consoleHelper.printColored("[!] IMPORTANT! Your beatmapcacheexpire in config.ini is > 0 and osu!api features are disabled.\nWe do not reccoment this, because too old beatmaps will be shown as unranked.\nSet beatmapcacheexpire to 0 to disable beatmap latest update check and fix that issue.", bcolors.YELLOW)

		# Load achievements
		try:
			secret.achievements.utils.load_achievements()
		except Exception as e:
			consoleHelper.printError()
			consoleHelper.printColored(
				"[!] Error while loading achievements! ({})".format(e),
				bcolors.RED,
			)
			sys.exit()
		consoleHelper.printDone()

		# Set achievements version
		glob.redis.set("lets:achievements_version", glob.ACHIEVEMENTS_VERSION)
		consoleHelper.printColored("Achievements version {}".format(glob.ACHIEVEMENTS_VERSION), bcolors.YELLOW)

		# Print disallowed mods into console (Used to also assign it into variable but has been moved elsewhere)
		unranked_mods = [key for key, value in glob.conf.extra["common"]["rankable-mods"].items() if not value]
		consoleHelper.printColored("Unranked mods: {}".format(", ".join(unranked_mods)), bcolors.YELLOW)
		
		# Print allowed beatmap rank statuses
		allowed_beatmap_rank = [key for key, value in glob.conf.extra["lets"]["allowed-beatmap-rankstatus"].items() if value]
		consoleHelper.printColored("Allowed beatmap rank statuses: {}".format(", ".join(allowed_beatmap_rank)), bcolors.YELLOW)

		# Make array of bools to respective rank id's
		glob.conf.extra["_allowed_beatmap_rank"] = [getattr(rankedStatuses, key) for key in allowed_beatmap_rank] # Store the allowed beatmap rank id's into glob


		# Discord
		if generalUtils.stringToBool(glob.conf.config["discord"]["enable"]):
			glob.schiavo = schiavo.schiavo(glob.conf.config["discord"]["boturl"], "**lets**")
		else:
			consoleHelper.printColored("[!] Warning! Discord logging is disabled!", bcolors.YELLOW)

		# Check debug mods
		glob.debug = generalUtils.stringToBool(glob.conf.config["server"]["debug"])
		if glob.debug:
			consoleHelper.printColored("[!] Warning! Server running in debug mode!", bcolors.YELLOW)

		# Server port
		try:
			serverPort = int(glob.conf.config["server"]["port"])
		except:
			consoleHelper.printColored("[!] Invalid server port! Please check your config.ini and run the server again", bcolors.RED)

		# Make app
		glob.application = make_app()

		# Set up sentry
		try:
			glob.sentry = generalUtils.stringToBool(glob.conf.config["sentry"]["enable"])
			if glob.sentry:
				glob.application.sentry_client = AsyncSentryClient(glob.conf.config["sentry"]["dsn"], release=glob.VERSION)
			else:
				consoleHelper.printColored("[!] Warning! Sentry logging is disabled!", bcolors.YELLOW)
		except:
			consoleHelper.printColored("[!] Error while starting Sentry client! Please check your config.ini and run the server again", bcolors.RED)

		# Set up Datadog
		try:
			if generalUtils.stringToBool(glob.conf.config["datadog"]["enable"]):
				glob.dog = datadogClient.datadogClient(glob.conf.config["datadog"]["apikey"], glob.conf.config["datadog"]["appkey"])
			else:
				consoleHelper.printColored("[!] Warning! Datadog stats tracking is disabled!", bcolors.YELLOW)
		except:
			consoleHelper.printColored("[!] Error while starting Datadog client! Please check your config.ini and run the server again", bcolors.RED)

		# Connect to pubsub channels
		pubSub.listener(glob.redis, {
			"lets:beatmap_updates": beatmapUpdateHandler.handler(),
		}).start()

		# Server start message and console output
		consoleHelper.printColored("> L.E.T.S. is listening for clients on {}:{}...".format(glob.conf.config["server"]["host"], serverPort), bcolors.GREEN)
		log.logMessage("Server started!", discord="bunker", stdout=False)

		# Start Tornado
		glob.application.listen(serverPort, address=glob.conf.config["server"]["host"])
		tornado.ioloop.IOLoop.instance().start()
	finally:
		# Perform some clean up
		print("> Disposing server... ")
		glob.fileBuffers.flushAll()
		consoleHelper.printColored("Goodbye!", bcolors.GREEN)
