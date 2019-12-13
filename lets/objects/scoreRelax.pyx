import time

from objects import beatmap
from common import generalUtils
from common.constants import gameModes
from common.log import logUtils as log
from common.ripple import userUtils
from constants import rankedStatuses
from common.ripple import scoreUtils
from objects import glob
from pp import relaxoppai
from pp import rippoppai
from pp import wifipiano2
from pp import cicciobello


class score:
	PP_CALCULATORS = {
		gameModes.STD: relaxoppai.oppai,
		gameModes.TAIKO: rippoppai.oppai,
		gameModes.CTB: cicciobello.Cicciobello,
		gameModes.MANIA: wifipiano2.piano
	}
	__slots__ = ["scoreID", "playerName", "score", "maxCombo", "c50", "c100", "c300", "cMiss", "cKatu", "cGeki",
	             "fullCombo", "mods", "playerUserID","rank","date", "hasReplay", "fileMd5", "passed", "playDateTime",
	             "gameMode", "completed", "accuracy", "pp", "oldPersonalBest", "rankedScoreIncrease", "personalOldBestScore"]
	def __init__(self, scoreID = None, rank = None, setData = True):
		"""
		Initialize a (empty) score object.

		scoreID -- score ID, used to get score data from db. Optional.
		rank -- score rank. Optional
		setData -- if True, set score data from db using scoreID. Optional.
		"""
		self.scoreID = 0
		self.playerName = "nospe"
		self.score = 0
		self.maxCombo = 0
		self.c50 = 0
		self.c100 = 0
		self.c300 = 0
		self.cMiss = 0
		self.cKatu = 0
		self.cGeki = 0
		self.fullCombo = False
		self.mods = 0
		self.playerUserID = 0
		self.rank = rank	# can be empty string too
		self.date = 0
		self.hasReplay = 0

		self.fileMd5 = None
		self.passed = False
		self.playDateTime = 0
		self.gameMode = 0
		self.completed = 0

		self.accuracy = 0.00

		self.pp = 0.00

		self.oldPersonalBest = 0
		self.rankedScoreIncrease = 0
		self.personalOldBestScore = None

		if scoreID is not None and setData:
			self.setDataFromDB(scoreID, rank)

	def calculateAccuracy(self):
		"""
		Calculate and set accuracy for that score
		"""
		if self.gameMode == 0:
			# std
			totalPoints = self.c50*50+self.c100*100+self.c300*300
			totalHits = self.c300+self.c100+self.c50+self.cMiss
			if totalHits == 0:
				self.accuracy = 1
			else:
				self.accuracy = totalPoints/(totalHits*300)
		elif self.gameMode == 1:
			# taiko
			totalPoints = (self.c100*50)+(self.c300*100)
			totalHits = self.cMiss+self.c100+self.c300
			if totalHits == 0:
				self.accuracy = 1
			else:
				self.accuracy = totalPoints / (totalHits * 100)
		elif self.gameMode == 2:
			# ctb
			fruits = self.c300+self.c100+self.c50
			totalFruits = fruits+self.cMiss+self.cKatu
			if totalFruits == 0:
				self.accuracy = 1
			else:
				self.accuracy = fruits / totalFruits
		elif self.gameMode == 3:
			# mania
			totalPoints = self.c50*50+self.c100*100+self.cKatu*200+self.c300*300+self.cGeki*300
			totalHits = self.cMiss+self.c50+self.c100+self.c300+self.cGeki+self.cKatu
			self.accuracy = totalPoints / (totalHits * 300)
		else:
			# unknown gamemode
			self.accuracy = 0

	def setRank(self, rank):
		"""
		Force a score rank

		rank -- new score rank
		"""
		self.rank = rank
			
	def setDataFromDB(self, scoreID, rank = None):
		"""
		Set this object's score data from db
		Sets playerUserID too

		scoreID -- score ID
		rank -- rank in scoreboard. Optional.
		"""
		data = glob.db.fetch("SELECT scores_relax.*, users.username FROM scores_relax LEFT JOIN users ON users.id = scores_relax.userid WHERE scores_relax.id = %s LIMIT 1", [scoreID])
		if data is not None:
			self.setDataFromDict(data, rank)

	def setDataFromDict(self, data, rank = None):
		"""
		Set this object's score data from dictionary
		Doesn't set playerUserID

		data -- score dictionarty
		rank -- rank in scoreboard. Optional.
		"""
		#print(str(data))
		self.scoreID = data["id"]
		if "username" in data:
			self.playerName = userUtils.getClan(data["userid"])
		else:
			self.playerName = userUtils.getUsername(data["userid"])
		self.playerUserID = data["userid"]
		self.score = data["score"]
		self.maxCombo = data["max_combo"]
		self.gameMode = data["play_mode"]
		self.c50 = data["50_count"]
		self.c100 = data["100_count"]
		self.c300 = data["300_count"]
		self.cMiss = data["misses_count"]
		self.cKatu = data["katus_count"]
		self.cGeki = data["gekis_count"]
		self.fullCombo = True if data["full_combo"] == 1 else False
		self.mods = data["mods"]
		self.rank = rank if rank is not None else ""
		self.date = data["time"]
		self.fileMd5 = data["beatmap_md5"]
		self.completed = data["completed"]
		#if "pp" in data:
		self.pp = data["pp"]
		self.calculateAccuracy()

	def setDataFromScoreData(self, scoreData):
		"""
		Set this object's score data from scoreData list (submit modular)

		scoreData -- scoreData list
		"""
		if len(scoreData) >= 16:
			self.fileMd5 = scoreData[0]
			self.playerName = scoreData[1].strip()
			# %s%s%s = scoreData[2]
			self.c300 = int(scoreData[3])
			self.c100 = int(scoreData[4])
			self.c50 = int(scoreData[5])
			self.cGeki = int(scoreData[6])
			self.cKatu = int(scoreData[7])
			self.cMiss = int(scoreData[8])
			self.score = int(scoreData[9])
			self.maxCombo = int(scoreData[10])
			self.fullCombo = True if scoreData[11] == 'True' else False
			#self.rank = scoreData[12]
			self.mods = int(scoreData[13])
			self.passed = True if scoreData[14] == 'True' else False
			self.gameMode = int(scoreData[15])
			#self.playDateTime = int(scoreData[16])
			self.playDateTime = int(time.time())
			self.calculateAccuracy()
			#osuVersion = scoreData[17]
			self.calculatePP()

			# Set completed status
			self.setCompletedStatus()


	def getData(self, pp=True):
		"""Return score row relative to this score for getscores"""
		return "{}|{}|{}|{}|{}|{}|{}|{}|{}|{}|{}|{}|{}|{}|{}|1\n".format(
			self.scoreID,
			self.playerName,
			int(self.pp) if pp else self.score,
			self.maxCombo,
			self.c50,
			self.c100,
			self.c300,
			self.cMiss,
			self.cKatu,
			self.cGeki,
			self.fullCombo,
			self.mods,
			self.playerUserID,
			self.rank,
			self.date)

	def setCompletedStatus(self):
		"""
		Set this score completed status and rankedScoreIncrease
		"""
		self.completed = 0
		if self.passed and scoreUtils.isRankable(self.mods):
			# Get userID
			userID = userUtils.getID(self.playerName)

			# Make sure we don't have another score identical to this one
			duplicate = glob.db.fetch("SELECT id FROM scores_relax WHERE userid = %s AND beatmap_md5 = %s AND play_mode = %s AND time = %s AND score = %s LIMIT 1", [userID, self.fileMd5, self.gameMode, self.date, self.score])
			if duplicate is not None:
				# Found same score in db. Don't save this score.
				self.completed = -1
				return

			# No duplicates found.
			# Get right "completed" value
			personalBest = glob.db.fetch("SELECT id,{}score FROM scores WHERE userid = %s AND beatmap_md5 = %s AND play_mode = %s AND completed = 3 LIMIT 1".format(
					glob.conf.extra["lets"]["submit"]["score-overwrite"] == "score" and " " or " {}, ".format(glob.conf.extra["lets"]["submit"]["score-overwrite"])
				),
				[userID, self.fileMd5, self.gameMode])
			if personalBest is None:
				# This is our first score on this map, so it's our best score
				self.completed = 3
				self.rankedScoreIncrease = self.score
				self.oldPersonalBest = 0
				self.personalOldBestScore = None
			else:
				self.personalOldBestScore = personalBest["id"]
				self.calculatePP()
				# Compare personal best's score with current score
				if getattr(self, glob.conf.extra["lets"]["submit"]["score-overwrite"]) > personalBest[glob.conf.extra["lets"]["submit"]["score-overwrite"]]:
					# New best score
					self.completed = 3
					self.rankedScoreIncrease = self.score-personalBest["score"]
					self.oldPersonalBest = personalBest["id"]
				else:
					self.completed = 2
					self.rankedScoreIncrease = 0
					self.oldPersonalBest = 0
				
		log.debug("Completed status: {}".format(self.completed))

	def saveScoreInDB(self):
		"""
		Save this score in DB (if passed and mods are valid)
		"""
		# Add this score
		if self.completed >= 2:
			query = "INSERT INTO scores_relax (id, beatmap_md5, userid, score, max_combo, full_combo, mods, 300_count, 100_count, 50_count, katus_count, gekis_count, misses_count, time, play_mode, completed, accuracy, pp) VALUES (NULL, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);"
			self.scoreID = int(glob.db.execute(query, [self.fileMd5, userUtils.getID(self.playerName), self.score, self.maxCombo, int(self.fullCombo), self.mods, self.c300, self.c100, self.c50, self.cKatu, self.cGeki, self.cMiss, self.playDateTime, self.gameMode, self.completed, self.accuracy * 100, self.pp]))

			# Set old personal best to completed = 2
			if self.oldPersonalBest != 0 and self.completed == 3:
				glob.db.execute("UPDATE scores_relax SET completed = 2 WHERE id = %s LIMIT 1", [self.oldPersonalBest])

	def calculatePP(self, b = None):
		"""
		Calculate this score's pp value if completed == 3
		"""
		# Create beatmap object
		if b is None:
			b = beatmap.beatmap(self.fileMd5, 0)

		# Calculate pp
		if b.rankedStatus >= rankedStatuses.RANKED and b.rankedStatus != rankedStatuses.UNKNOWN \
			and scoreUtils.isRankable(self.mods) and self.passed and self.gameMode in score.PP_CALCULATORS:
			calculator = score.PP_CALCULATORS[self.gameMode](b, self)
			self.pp = calculator.pp
		else:
			self.pp = 0
