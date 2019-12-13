import tornado.gen
import tornado.web

from common.log import logUtils as log
from common.web import requestsManager
from constants import exceptions
from helpers import osuapiHelper
from common.sentry import sentry

MODULE_NAME = "maps"
class handler(requestsManager.asyncRequestHandler):
	@tornado.web.asynchronous
	@tornado.gen.engine
	@sentry.captureTornado
	def asyncGet(self, fileName = None):
		try:
			# Check arguments
			if fileName is None:
				raise exceptions.invalidArgumentsException(MODULE_NAME)
			if fileName == "":
				raise exceptions.invalidArgumentsException(MODULE_NAME)

			fileNameShort = fileName[:32]+"..." if len(fileName) > 32 else fileName[:-4]
			log.info("Requested .osu file {}".format(fileNameShort))

			# Get .osu file from osu! server
			fileContent = osuapiHelper.getOsuFileFromName(fileName)
			if fileContent is None:
				# TODO: Sentry capture message here
				raise exceptions.osuApiFailException(MODULE_NAME)
			self.write(fileContent)
		except exceptions.invalidArgumentsException:
			self.set_status(500)
		except exceptions.osuApiFailException:
			self.set_status(500)