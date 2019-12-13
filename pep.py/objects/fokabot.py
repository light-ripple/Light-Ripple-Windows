"""FokaBot related functions"""
import re

from common import generalUtils
from common.constants import actions
from common.ripple import userUtils
from constants import fokabotCommands
from constants import serverPackets
from objects import glob

# Tillerino np regex, compiled only once to increase performance
npRegex = re.compile("^https?:\\/\\/osu\\.ppy\\.sh\\/b\\/(\\d*)")

def connect(timeOffset = 9):
	"""
	Connect FokaBot to Bancho

	:return:
	"""
	glob.BOT_NAME = userUtils.getUsername(999)
	token = glob.tokens.addToken(999)
	token.actionID = actions.IDLE
	token.actionText = "\n-- Welcome to Ainu --"
	token.pp = 727
	token.accuracy = 0.9885
	token.playcount = 26956
	token.totalScore = 237228316533
	token.timeOffset = timeOffset
	token.timezone = 24+token.timeOffset
	token.country = 111
	glob.streams.broadcast("main", serverPackets.userPanel(999))
	glob.streams.broadcast("main", serverPackets.userStats(999))

def disconnect():
	"""
	Disconnect FokaBot from Bancho

	:return:
	"""
	glob.tokens.deleteToken(glob.tokens.getTokenFromUserID(999))

def fokabotResponse(fro, chan, message):
	"""
	Check if a message has triggered FokaBot

	:param fro: sender username
	:param chan: channel name (or receiver username)
	:param message: chat mesage
	:return: FokaBot's response or False if no response
	"""
	for i in fokabotCommands.commands:
		# Loop though all commands
		if re.compile("^{}( (.+)?)?$".format(i["trigger"])).match(message.strip()):
			# message has triggered a command

			# Make sure the user has right permissions
			if i["privileges"] is not None:
				# Rank = x
				if userUtils.getPrivileges(userUtils.getID(fro)) & i["privileges"] == 0:
					return False

			# Check argument number
			message = message.split(" ")
			if i["syntax"] != "" and len(message) <= len(i["syntax"].split(" ")):
				return "Wrong syntax: {} {}".format(i["trigger"], i["syntax"])

			# Return response or execute callback
			if i["callback"] is None:
				return i["response"]
			else:
				return i["callback"](fro, chan, message[1:])

	# No commands triggered
	return False