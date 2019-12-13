from objects import score
from common.ripple import userUtils
from constants import rankedStatuses
from common.constants import mods as modsEnum
from objects import glob


class scoreboard:
	def __init__(self, username, gameMode, beatmap, setScores = True, country = False, friends = False, mods = -1):
		"""
		Initialize a leaderboard object

		username -- username of who's requesting the scoreboard. None if not known
		gameMode -- requested gameMode
		beatmap -- beatmap objecy relative to this leaderboard
		setScores -- if True, will get personal/top 50 scores automatically. Optional. Default: True
		"""
		self.scores = []				# list containing all top 50 scores objects. First object is personal best
		self.totalScores = 0
		self.personalBestRank = -1		# our personal best rank, -1 if not found yet
		self.username = username		# username of who's requesting the scoreboard. None if not known
		self.userID = userUtils.getID(self.username)	# username's userID
		self.gameMode = gameMode		# requested gameMode
		self.beatmap = beatmap			# beatmap objecy relative to this leaderboard
		self.country = country
		self.friends = friends
		self.mods = mods
		if setScores:
			self.setScores()

	@staticmethod
	def buildQuery(params):
		return "{select} {joins} {country} {mods} {friends} {order} {limit}".format(**params)
		
	def getPersonalBest(self):
		if self.userID == 0:
			return None
			
		# Query parts
		cdef str select = ""
		cdef str joins = ""
		cdef str country = ""
		cdef str mods = ""
		cdef str friends = ""
		cdef str order = ""
		cdef str limit = ""
		select = "SELECT id FROM scores WHERE userid = %(userid)s AND beatmap_md5 = %(md5)s AND play_mode = %(mode)s AND completed = 3"

		# Mods
		if self.mods > -1:
			mods = "AND mods = %(mods)s"

		# Friends ranking
		if self.friends:
			friends = "AND (scores.userid IN (SELECT user2 FROM users_relationships WHERE user1 = %(userid)s) OR scores.userid = %(userid)s)"

		# Sort and limit at the end
		order = "ORDER BY score DESC"
		limit = "LIMIT 1"

		# Build query, get params and run query
		query = self.buildQuery(locals())
		params = {"userid": self.userID, "md5": self.beatmap.fileMD5, "mode": self.gameMode, "mods": self.mods}
		id_ = glob.db.fetch(query, params)
		if id_ is None:
			return None
		return id_["id"]
			
	def setScores(self):
		"""
		Set scores list
		"""
		
		# Reset score list
		self.scores = []
		self.scores.append(-1)

		# Make sure the beatmap is ranked
		if self.beatmap.rankedStatus not in glob.conf.extra["_allowed_beatmap_rank"]:
			return

		# Query parts
		cdef str select = ""
		cdef str joins = ""
		cdef str country = ""
		cdef str mods = ""
		cdef str friends = ""
		cdef str order = ""
		cdef str limit = ""

		# Find personal best score
		personalBestScore = self.getPersonalBest()

		# Output our personal best if found
		if personalBestScore is not None:
			s = score.score(personalBestScore)
			self.scores[0] = s
		else:
			# No personal best
			self.scores[0] = -1

		# Get top 50 scores
		select = "SELECT *"
		joins = "FROM scores STRAIGHT_JOIN users ON scores.userid = users.id STRAIGHT_JOIN users_stats ON users.id = users_stats.id WHERE scores.beatmap_md5 = %(beatmap_md5)s AND scores.play_mode = %(play_mode)s AND scores.completed = 3 AND (users.privileges & 1 > 0 OR users.id = %(userid)s)"

		# Country ranking
		if self.country:
			country = "AND users_stats.country = (SELECT country FROM users_stats WHERE id = %(userid)s LIMIT 1)"
		else:
			country = ""

		# Mods ranking (ignore auto, since we use it for pp sorting)
		if self.mods > -1 and self.mods & modsEnum.AUTOPLAY == 0:
			mods = "AND scores.mods = %(mods)s"
		else:
			mods = ""

		# Friends ranking
		if self.friends:
			friends = "AND (scores.userid IN (SELECT user2 FROM users_relationships WHERE user1 = %(userid)s) OR scores.userid = %(userid)s)"
		else:
			friends = ""

		# Sort and limit at the end
		if not glob.conf.extra["lets"]["scoreboard"]["ppboard"] and self.mods <= -1 or self.mods & modsEnum.AUTOPLAY == 0:
			# Order by score if we aren't filtering by mods or autoplay mod is disabled
			order = "ORDER BY score DESC"
		elif self.mods & modsEnum.AUTOPLAY > 0 or glob.conf.extra["lets"]["scoreboard"]["ppboard"]:
			# Otherwise, filter by pp
			order = "ORDER BY pp DESC"
		limit = "LIMIT 50"

		# Build query, get params and run query
		query = self.buildQuery(locals())
		params = {"beatmap_md5": self.beatmap.fileMD5, "play_mode": self.gameMode, "userid": self.userID, "mods": self.mods}
		topScores = glob.db.fetchAll(query, params)

		# Set data for all scores
		cdef int c = 1
		cdef dict topScore
		if topScores is not None:
			for topScore in topScores:
				# Create score object
				s = score.score(topScore["id"], setData=False)

				# Set data and rank from topScores's row
				s.setDataFromDict(topScore)
				s.rank = c

				# Check if this top 50 score is our personal best
				if s.playerName == self.username:
					self.personalBestRank = c

				# Add this score to scores list and increment rank
				self.scores.append(s)
				c+=1

		'''# If we have more than 50 scores, run query to get scores count
		if c >= 50:
			# Count all scores on this map
			select = "SELECT COUNT(*) AS count"
			limit = "LIMIT 1"

			# Build query, get params and run query
			query = self.buildQuery(locals())
			count = glob.db.fetch(query, params)
			if count == None:
				self.totalScores = 0
			else:
				self.totalScores = count["count"]
		else:
			self.totalScores = c-1'''

		# If personal best score was not in top 50, try to get it from cache
		if personalBestScore is not None and self.personalBestRank < 1:
			self.personalBestRank = glob.personalBestCache.get(self.userID, self.beatmap.fileMD5, self.country, self.friends, self.mods)

		# It's not even in cache, get it from db
		if personalBestScore is not None and self.personalBestRank < 1:
			self.setPersonalBest()

		# Cache our personal best rank so we can eventually use it later as
		# before personal best rank" in submit modular when building ranking panel
		if self.personalBestRank >= 1:
			glob.personalBestCache.set(self.userID, self.personalBestRank, self.beatmap.fileMD5)

	def setPersonalBest(self):
		"""
		Set personal best rank ONLY
		Ikr, that query is HUGE but xd
		"""
		# Before running the HUGE query, make sure we have a score on that map
		cdef str query = "SELECT id FROM scores WHERE beatmap_md5 = %(md5)s AND userid = %(userid)s AND play_mode = %(mode)s AND completed = 3"
		# Mods
		if self.mods > -1:
			query += " AND scores.mods = %(mods)s"
		# Friends ranking
		if self.friends:
			query += " AND (scores.userid IN (SELECT user2 FROM users_relationships WHERE user1 = %(userid)s) OR scores.userid = %(userid)s)"
		# Sort and limit at the end
		query += " LIMIT 1"
		hasScore = glob.db.fetch(query, {"md5": self.beatmap.fileMD5, "userid": self.userID, "mode": self.gameMode, "mods": self.mods})
		if hasScore is None:
			return

		overwrite = glob.conf.extra["lets"]["scoreboard"]["ppboard"] and "pp" or "score"
		
		# We have a score, run the huge query
		# Base query
		query = """SELECT COUNT(*) AS rank FROM scores STRAIGHT_JOIN users ON scores.userid = users.id STRAIGHT_JOIN users_stats ON users.id = users_stats.id WHERE scores.{0} >= (
		SELECT {0} FROM scores WHERE beatmap_md5 = %(md5)s AND play_mode = %(mode)s AND completed = 3 AND userid = %(userid)s LIMIT 1
		) AND scores.beatmap_md5 = %(md5)s AND scores.play_mode = %(mode)s AND scores.completed = 3 AND users.privileges & 1 > 0""".format(overwrite)
		# Country
		if self.country:
			query += " AND users_stats.country = (SELECT country FROM users_stats WHERE id = %(userid)s LIMIT 1)"
		# Mods
		if self.mods > -1:
			query += " AND scores.mods = %(mods)s"
		# Friends
		if self.friends:
			query += " AND (scores.userid IN (SELECT user2 FROM users_relationships WHERE user1 = %(userid)s) OR scores.userid = %(userid)s)"
		# Sort and limit at the end
		query += " ORDER BY {} DESC LIMIT 1".format(overwrite)
		result = glob.db.fetch(query, {"md5": self.beatmap.fileMD5, "userid": self.userID, "mode": self.gameMode, "mods": self.mods})
		if result is not None:
			self.personalBestRank = result["rank"]

	def getScoresData(self):
		"""
		Return scores data for getscores

		return -- score data in getscores format
		"""
		data = ""

		# Output personal best
		if self.scores[0] == -1:
			# We don't have a personal best score
			data += "\n"
		else:
			# Set personal best score rank
			self.setPersonalBest()	# sets self.personalBestRank with the huge query
			self.scores[0].setRank(self.personalBestRank)
			data += self.scores[0].getData(pp=glob.conf.extra["lets"]["scoreboard"]["ppboard"])

		# Output top 50 scores
		for i in self.scores[1:]:
			data += i.getData(pp=glob.conf.extra["lets"]["scoreboard"]["ppboard"] or (self.mods > -1 and self.mods & modsEnum.AUTOPLAY > 0))

		return data
