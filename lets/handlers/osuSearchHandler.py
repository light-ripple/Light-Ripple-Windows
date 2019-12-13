import tornado.gen
import tornado.web

from common.sentry import sentry
from common.web import requestsManager
from common.web import cheesegull
from constants import exceptions
from common.log import logUtils as log

MODULE_NAME = "direct"
class handler(requestsManager.asyncRequestHandler):
	"""
	Handler for /web/osu-search.php
	"""
	@tornado.web.asynchronous
	@tornado.gen.engine
	@sentry.captureTornado
	def asyncGet(self):
		output = ""
		try:
			try:
				# Get arguments
				gameMode = self.get_argument("m", None)
				if gameMode is not None:
					gameMode = int(gameMode)
				if gameMode < 0 or gameMode > 3:
					gameMode = None

				rankedStatus = self.get_argument("r", None)
				if rankedStatus is not None:
					rankedStatus = int(rankedStatus)

				query = self.get_argument("q", "")
				page = int(self.get_argument("p", "0"))
				if query.lower() in ["newest", "top rated", "most played"]:
					query = ""
			except ValueError:
				raise exceptions.invalidArgumentsException(MODULE_NAME)

			# Get data from cheesegull API
			log.info("Requested osu!direct search: {}".format(query if query != "" else "index"))
			searchData = cheesegull.getListing(rankedStatus=cheesegull.directToApiStatus(rankedStatus), page=page * 100, gameMode=gameMode, query=query)
			if searchData is None or searchData is None:
				raise exceptions.noAPIDataError()

			# Write output
			output += "999" if len(searchData) == 100 else str(len(searchData))
			output += "\n"
			for beatmapSet in searchData:
				try:
					output += cheesegull.toDirect(beatmapSet) + "\r\n"
				except ValueError:
					# Invalid cheesegull beatmap (empty beatmapset, cheesegull bug? See Sentry #LETS-00-32)
					pass
		except (exceptions.noAPIDataError, exceptions.invalidArgumentsException):
			output = "0\n"
		finally:
			self.write(output)
