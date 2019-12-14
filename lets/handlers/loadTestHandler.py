import tornado.gen
import tornado.web

from common.web import requestsManager
from objects import glob


class handler(requestsManager.asyncRequestHandler):
	@tornado.web.asynchronous
	@tornado.gen.engine
	def asyncGet(self):
		if not glob.debug:
			self.write("Nope")
			return
		glob.db.fetchAll("SELECT SQL_NO_CACHE * FROM beatmaps")
		glob.db.fetchAll("SELECT SQL_NO_CACHE * FROM users")
		glob.db.fetchAll("SELECT SQL_NO_CACHE * FROM scores")
		self.write("ibmd")
