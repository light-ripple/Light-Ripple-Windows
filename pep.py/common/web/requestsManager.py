import sys
import traceback

import tornado
import tornado.web
import tornado.gen
from tornado.ioloop import IOLoop
from objects import glob
from common.log import logUtils as log
from raven.contrib.tornado import SentryMixin

class asyncRequestHandler(SentryMixin, tornado.web.RequestHandler):
	"""
	Tornado asynchronous request handler
	create a class that extends this one (requestHelper.asyncRequestHandler)
	use asyncGet() and asyncPost() instead of get() and post().
	Done. I'm not kidding.
	"""
	@tornado.web.asynchronous
	@tornado.gen.engine
	def get(self, *args, **kwargs):
		try:
			yield tornado.gen.Task(runBackground, (self.asyncGet, tuple(args), dict(kwargs)))
		finally:
			if not self._finished:
				self.finish()

	@tornado.web.asynchronous
	@tornado.gen.engine
	def post(self, *args, **kwargs):
		try:
			yield tornado.gen.Task(runBackground, (self.asyncPost, tuple(args), dict(kwargs)))
		finally:
			if not self._finished:
				self.finish()

	def asyncGet(self, *args, **kwargs):
		self.send_error(405)

	def asyncPost(self, *args, **kwargs):
		self.send_error(405)

	def getRequestIP(self):
		"""
		Return CF-Connecting-IP (request IP when under cloudflare, you have to configure nginx to enable that)
		If that fails, return X-Forwarded-For (request IP when not under Cloudflare)
		if everything else fails, return remote IP

		:return: Client IP address
		"""
		if "CF-Connecting-IP" in self.request.headers:
			return self.request.headers.get("CF-Connecting-IP")
		elif "X-Forwarded-For" in self.request.headers:
			return self.request.headers.get("X-Forwarded-For")
		else:
			return self.request.remote_ip


def runBackground(data, callback):
	"""
	Run a function in the background.
	Used to handle multiple requests at the same time

	:param data: (func, args, kwargs)
	:param callback: function to call when `func` (data[0]) returns
	:return:
	"""
	func, args, kwargs = data
	def _callback(result):
		IOLoop.instance().add_callback(lambda: callback(result))
	glob.pool.apply_async(func, args, kwargs, _callback)
	glob.dog.increment(glob.DATADOG_PREFIX + ".incoming_requests")

def checkArguments(arguments, requiredArguments):
	"""
	Check that every requiredArguments elements are in arguments

	:param arguments: full argument list, from tornado
	:param requiredArguments: required arguments list
	:return: True if all arguments are passed, False if not
	"""
	for i in requiredArguments:
		if i not in arguments:
			return False
	return True

def printArguments(t):
	"""
	Print passed arguments, for debug purposes

	:param t: tornado object (self)
	"""
	msg = "ARGS::"
	for i in t.request.arguments:
		msg += "{}={}\r\n".format(i, t.get_argument(i))
	log.debug(msg)
