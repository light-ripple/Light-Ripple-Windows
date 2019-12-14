import threading
import datadog
from objects import glob

class periodicCheck:
	def __init__(self, name, checkFunction):
		"""
		Initialize a periodic check object

		:param name: Datadog stat name, without prefix
		:param checkFunction: Function that returns the data to report. Eg: `lambda: len(something)`
		"""
		self.name = glob.DATADOG_PREFIX+"."+name
		self.checkFunction = checkFunction

class datadogClient:
	def __init__(self, apiKey=None, appKey=None, periodicChecks=None):
		"""
		Initialize a toggleable Datadog Client

		:param apiKey: Datadog api key. Leave empty to create a dummy (disabled) Datadog client.
		:param appKey: Datadog app key. Leave empty to create a dummy (disabled) Datadog client.
		:param periodicChecks: List of periodicCheck objects. Optional. Leave empty to disable periodic checks.
		"""
		if apiKey is not None and appKey is not None:
			datadog.initialize(api_key=apiKey, app_key=appKey)
			self.client = datadog.ThreadStats()
			self.client.start()
			self.periodicChecks = periodicChecks
			if self.periodicChecks is not None:
				threading.Thread(target=self.__periodicCheckLoop).start()
		else:
			self.client = None

	def increment(self, *args, **kwargs):
		"""
		Call self.client.increment(*args, **kwargs) if this client is not a dummy

		:param args:
		:param kwargs:
		:return:
		"""
		if self.client is not None:
			self.client.increment(*args, **kwargs)

	def gauge(self, *args, **kwargs):
		"""
		Call self.client.gauge(*args, **kwargs) if this client is not a dummy

		:param args:
		:param kwargs:
		:return:
		"""
		if self.client is not None:
			self.client.gauge(*args, **kwargs)

	def __periodicCheckLoop(self):
		"""
		Report periodic data to datadog.
		Called every 5 seconds.
		Call this function only once.

		:return:
		"""
		if self.periodicChecks is None:
			return

		# Get data
		for i in self.periodicChecks:
			self.gauge(i.name, i.checkFunction())

		# Schedule a new datadog update
		threading.Timer(10, self.__periodicCheckLoop).start()