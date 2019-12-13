from common.log import logUtils as log
from common.ripple import userUtils
from objects import glob
import json

class userStatsCache:
	def get(self, userID, gameMode):
		"""
		Get cached user stats from redis.
		If user stats are not cached, they'll be read from db, cached and returned

		:param userID: userID
		:param gameMode: game mode number
		:return: userStats dictionary (rankedScore, totalScore, pp, accuracy, playcount)
		"""
		data = glob.redis.get("lets:user_stats_cache:{}:{}".format(gameMode, userID))
		if data is None:
			# If data is not cached, cache it and call get function again
			log.debug("userStatsCache miss")
			self.update(userID, gameMode)
			return self.get(userID, gameMode)

		log.debug("userStatsCache hit")
		retData = json.loads(data.decode("utf-8"))
		return retData

	def update(self, userID, gameMode, data = None):
		"""
		Update cached user stats in redis with new values

		:param userID: userID
		:param gameMode: game mode number
		:param data: data to cache. Optional. If not passed, will get from db
		:return:
		"""
		if data is None:
			data = {}
		if len(data) == 0:
			data = userUtils.getUserStats(userID, gameMode)
		log.debug("userStatsCache set {}".format(data))
		glob.redis.set("lets:user_stats_cache:{}:{}".format(gameMode, userID), json.dumps(data), 1800)