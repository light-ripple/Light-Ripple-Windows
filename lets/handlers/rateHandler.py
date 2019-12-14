import tornado.gen
import tornado.web

from common.ripple import userUtils
from common.sentry import sentry
from common.web import requestsManager
from common.log import logUtils as log
from constants import exceptions, rankedStatuses
from objects import glob

MODULE_NAME = "rate"


class handler(requestsManager.asyncRequestHandler):
	@tornado.web.asynchronous
	@tornado.gen.engine
	@sentry.captureTornado
	def asyncGet(self):
		output = ""

		try:
			if not requestsManager.checkArguments(self.request.arguments, ["c", "u", "p"]):
				raise exceptions.invalidArgumentsException(MODULE_NAME)

			ip = self.getRequestIP()
			username = self.get_argument("u").strip()
			password = self.get_argument("p").strip()
			user_id = userUtils.getID(username)
			checksum = self.get_argument("c").strip()
			if not user_id:
				raise exceptions.loginFailedException(MODULE_NAME, user_id)
			if not userUtils.checkLogin(user_id, password, ip):
				raise exceptions.loginFailedException(MODULE_NAME, username)
			if userUtils.check2FA(user_id, ip):
				raise exceptions.need2FAException(MODULE_NAME, user_id, ip)

			ranked = glob.db.fetch(
				"SELECT ranked FROM beatmaps WHERE beatmap_md5 = %s LIMIT 1",
				(checksum,)
			)
			if ranked is None:
				output = "no exist"
				return
			if ranked["ranked"] < rankedStatuses.RANKED:
				output = "not ranked"
				return

			rating = glob.db.fetch("SELECT rating FROM beatmaps WHERE beatmap_md5 = %s LIMIT 1", (checksum,))
			has_voted = glob.db.fetch(
				"SELECT id FROM beatmaps_rating WHERE user_id = %s AND beatmap_md5 = %s LIMIT 1",
				(user_id, checksum)
			)
			if has_voted is not None:
				output = f"alreadyvoted\n{rating['rating']:.2f}"
				return
			vote = self.get_argument("v", default=None)
			if vote is None:
				output = "ok"
				return
			try:
				vote = int(vote)
			except ValueError:
				raise exceptions.invalidArgumentsException(MODULE_NAME)
			if vote < 0 or vote > 10:
				output = "out of range"
				return
			glob.db.execute(
				"REPLACE INTO beatmaps_rating (beatmap_md5, user_id, rating) VALUES (%s, %s, %s)",
				(checksum, user_id, vote)
			)
			glob.db.execute(
				"UPDATE beatmaps SET rating = (SELECT SUM(rating)/COUNT(rating) FROM beatmaps_rating "
				"WHERE beatmap_md5 = %(md5)s) WHERE beatmap_md5 = %(md5)s LIMIT 1",
				{"md5": checksum}
			)
			rating = glob.db.fetch("SELECT rating FROM beatmaps WHERE beatmap_md5 = %s LIMIT 1", (checksum,))
			output = f"{rating['rating']:.2f}"
		except exceptions.loginFailedException:
			output = "auth failed"
		except exceptions.invalidArgumentsException:
			output = "no"
		finally:
			self.write(output)