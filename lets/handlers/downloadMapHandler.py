import tornado.gen
import tornado.web

from common.web import requestsManager
from common.sentry import sentry

MODULE_NAME = "direct_download"
class handler(requestsManager.asyncRequestHandler):
	"""
	Handler for /d/
	"""
	@tornado.web.asynchronous
	@tornado.gen.engine
	@sentry.captureTornado
	def asyncGet(self, bid):
		try:
			noVideo = bid.endswith("n")
			if noVideo:
				bid = bid[:-1]
			bid = int(bid)

			self.set_status(302, "Moved Temporarily")
			url = "https://bm6.ppy.sh/d/{}{}".format(bid, "?novideo" if noVideo else "")
			self.add_header("Location", url)
			self.add_header("Cache-Control", "no-cache")
			self.add_header("Pragma", "no-cache")
		except ValueError:
			self.set_status(400)
			self.write("Invalid set id")