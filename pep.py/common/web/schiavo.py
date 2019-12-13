import requests
from urllib.parse import urlencode

class schiavo:
	"""
	Schiavo Bot class
	"""
	def __init__(self, botURL=None, prefix="", maxRetries=20):
		"""
		Initialize a new schiavo bot instance

		:param botURL: schiavo api url.
		:param prefix: text to prepend in every message, can be empty.
		:param maxRetries: max retries if api request fail. 0 = don't retry.
		"""
		self.botURL = botURL
		self.maxRetries = maxRetries
		self.prefix = prefix

	def sendMessage(self, channel, message, noPrefix=False):
		"""
		Send a generic message through schiavo api

		:param channel: api channel.
		:param message: message content.
		:param noPrefix: if True, don't prepend prefix to message.
		:return:
		"""
		if self.botURL is None:
			return
		for _ in range(0, self.maxRetries):
			try:
				finalMsg = "{prefix} {message}".format(prefix=self.prefix if not noPrefix else "", message=message)
				requests.get("{}/{}?{}".format(self.botURL, channel, urlencode({ "message": finalMsg })))
				break
			except requests.RequestException:
				continue

	def sendConfidential(self, message, noPrefix=False):
		"""
		Send a message to #bunk

		:param message: message content.
		:param noPrefix: if True, don't prepend prefix to message.
		:return:
		"""
		self.sendMessage("bunk", message, noPrefix)

	def sendStaff(self, message, noPrefix=False):
		"""
		Send a message to #staff

		:param message: message content.
		:param noPrefix: if True, don't prepend prefix to message.
		:return:
		"""
		self.sendMessage("staff", message, noPrefix)

	def sendGeneral(self, message, noPrefix=True):
		"""
		Send a message to #general

		:param message: message content.
		:param noPrefix: if True, don't prepend prefix to message.
		:return:
		"""
		self.sendMessage("general", message, noPrefix)

	def sendChatlog(self, message, noPrefix=True):
		"""
		Send a message to #chatlog.

		:param message: message content.
		:param noPrefix: if True, don't prepend prefix to message.
		:return:
		"""
		self.sendMessage("chatlog", message, noPrefix)

	def sendCM(self, message, noPrefix=False):
		"""
		Send a message to #communitymanagers

		:param message: message content.
		:param noPrefix: if True, don't prepend prefix to message.
		:return:
		"""
		self.sendMessage("cm", message, noPrefix)