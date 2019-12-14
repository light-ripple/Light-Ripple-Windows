import threading
from common.log import logUtils as log
from common.redis import generalPubSubHandler
from common.sentry import sentry

class listener(threading.Thread):
	def __init__(self, r, handlers):
		"""
		Initialize a set of redis pubSub listeners

		:param r: redis instance (usually glob.redis)
		:param handlers: dictionary with the following structure:
		```
		{
			"redis_channel_name": handler,
			...
		}
		```
		Where handler is:
		- 	An object of a class that inherits common.redis.generalPubSubHandler.
			You can create custom behaviors for your handlers by overwriting the `handle(self, data)` method,
			that will be called when that handler receives some data.

		- 	A function *object (not call)* that accepts one argument, that'll be the data received through the channel.
			This is useful if you want to make some simple handlers through a lambda, without having to create a class.
		"""
		threading.Thread.__init__(self)
		self.redis = r
		self.pubSub = self.redis.pubsub()
		self.handlers = handlers
		channels = []
		for k, v in self.handlers.items():
			channels.append(k)
		self.pubSub.subscribe(channels)
		log.info("Subscribed to redis pubsub channels: {}".format(channels))

	@sentry.capture()
	def processItem(self, item):
		"""
		Processes a pubSub item by calling channel's handler

		:param item: incoming data
		:return:
		"""
		if item["type"] == "message":
			# Process the message only if the channel has received a message
			# Decode the message
			item["channel"] = item["channel"].decode("utf-8")

			# Make sure the handler exists
			if item["channel"] in self.handlers:
				log.info("Redis pubsub: {} <- {} ".format(item["channel"], item["data"]))
				if isinstance(self.handlers[item["channel"]], generalPubSubHandler.generalPubSubHandler):
					# Handler class
					self.handlers[item["channel"]].handle(item["data"])
				else:
					# Function
					self.handlers[item["channel"]](item["data"])

	def run(self):
		"""
		Listen for data on incoming channels and process it.
		Runs forever.

		:return:
		"""
		for item in self.pubSub.listen():
			self.processItem(item)