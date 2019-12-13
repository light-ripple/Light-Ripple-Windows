from urllib.parse import urlencode

import requests
import tornado.gen
import tornado.web

from common.log import logUtils as log
from common.web import requestsManager


class handler(requestsManager.asyncRequestHandler):
	@tornado.web.asynchronous
	@tornado.gen.engine
	def asyncGet(self):
		try:
			args = {}
			#if "stream" in self.request.arguments:
			#	args["stream"] = self.get_argument("stream")
			#if "action" in self.request.arguments:
			#	args["action"] = self.get_argument("action")
			#if "time" in self.request.arguments:
			#	args["time"] = self.get_argument("time")

			# Pass all arguments otherwise it doesn't work
			for key, _ in self.request.arguments.items():
				args[key] = self.get_argument(key)

			if args["action"].lower() == "put":
				self.write("nope")
				return

			response = requests.get("https://osu.ppy.sh/web/check-updates.php?{}".format(urlencode(args)))
			self.write(response.text)
		except Exception as e:
			log.error("check-updates failed: {}".format(e))
			self.write("")
