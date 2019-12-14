from common.log import logUtils as log


class invalidArgumentsException(Exception):
	def __init__(self, handler):
		log.warning("{} - Invalid arguments".format(handler))

class loginFailedException(Exception):
	def __init__(self, handler, who):
		log.warning("{} - {}'s Login failed".format(handler, who))

class userBannedException(Exception):
	def __init__(self, handler, who):
		log.warning("{} - {} is banned".format(handler, who))

class userLockedException(Exception):
	def __init__(self, handler, who):
		log.warning("{} - {} is locked".format(handler, who))

class noBanchoSessionException(Exception):
	def __init__(self, handler, who, ip):
		log.warning("{handler} - {username} has tried to submit a score from {ip} without an active bancho session from that ip. If this happens often, {username} is trying to use a score submitter.".format(handler=handler, ip=ip, username=who), "bunker")

class osuApiFailException(Exception):
	def __init__(self, handler):
		log.warning("{} - Invalid data from osu!api".format(handler))

class fileNotFoundException(Exception):
	def __init__(self, handler, f):
		log.warning("{} - File not found ({})".format(handler, f))

class invalidBeatmapException(Exception):
	pass

class unsupportedGameModeException(Exception):
	pass

class beatmapTooLongException(Exception):
	def __init__(self, handler):
		log.warning("{} - Requested beatmap is too long.".format(handler))

class need2FAException(Exception):
	def __init__(self, handler, who, ip):
		log.warning("{} - 2FA check needed for user {} ({})".format(handler, who, ip))

class noAPIDataError(Exception):
	pass

class scoreNotFoundError(Exception):
	pass

class ppCalcException(Exception):
	def __init__(self, exception):
		self.exception = exception