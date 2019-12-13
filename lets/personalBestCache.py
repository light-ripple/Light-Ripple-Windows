from common.log import logUtils as log
from common import generalUtils
from objects import glob

class cacheMiss(Exception):
	pass

class personalBestCache:
	def get(self, userID, fileMd5, country=False, friends=False, mods=-1):
		"""
		Get cached personal best rank

		:param userID: userID
		:param fileMd5: beatmap md5
		:param country: True if country leaderboard, otherwise False
		:param friends: True if friends leaderboard, otherwise False
		:param mods: leaderboard mods
		:return: 0 if cache miss, otherwise rank number
		"""
		try:
			# Make sure the value is in cache
			data = glob.redis.get("lets:personal_best_cache:{}".format(userID))
			if data is None:
				raise cacheMiss()

			# Unpack cached data
			data = data.decode("utf-8").split("|")
			cachedpersonalBestRank = int(data[0])
			cachedfileMd5 = str(data[1])
			cachedCountry = generalUtils.stringToBool(data[2])
			cachedFriends = generalUtils.stringToBool(data[3])
			cachedMods = int(data[4])

			# Check if everything matches
			if fileMd5 != cachedfileMd5 or country != cachedCountry or friends != cachedFriends or mods != cachedMods:
				raise cacheMiss()

			# Cache hit
			log.debug("personalBestCache hit")
			return cachedpersonalBestRank
		except cacheMiss:
			log.debug("personalBestCache miss")
			return 0

	def set(self, userID, rank, fileMd5, country=False, friends=False, mods=-1):
		"""
		Set userID's redis personal best cache

		:param userID: userID
		:param rank: leaderboard rank
		:param fileMd5: beatmap md5
		:param country: True if country leaderboard, otherwise False
		:param friends: True if friends leaderboard, otherwise False
		:param mods: leaderboard mods
		:return:
		"""
		glob.redis.set("lets:personal_best_cache:{}".format(userID), "{}|{}|{}|{}|{}".format(rank, fileMd5, country, friends, mods), 1800)
		log.debug("personalBestCache set")
