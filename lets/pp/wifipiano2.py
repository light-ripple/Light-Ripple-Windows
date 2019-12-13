"""
Wifipiano 2

This file has been written taking by reference code from
osu-performance (https://github.com/ppy/osu-performance)
by Tom94, licensed under the GNU AGPL 3 License.
"""
from common.constants import mods
from common.log import logUtils as log
from constants import exceptions
from helpers import mapsHelper


class piano:
	__slots__ = ["beatmap", "score", "pp"]

	def __init__(self, __beatmap, __score):
		self.beatmap = __beatmap
		self.score = __score
		self.pp = 0
		self.getPP()

	def getPP(self):
		try:
			stars = self.beatmap.starsMania
			if stars == 0:
				# This beatmap can't be converted to mania
				raise exceptions.invalidBeatmapException()

			# Cache beatmap for cono
			mapFile = mapsHelper.cachedMapPath(self.beatmap.beatmapID)
			mapsHelper.cacheMap(mapFile, self.beatmap)

			od = self.beatmap.OD
			objects = self.score.c50+self.score.c100+self.score.c300+self.score.cKatu+self.score.cGeki+self.score.cMiss

			score = self.score.score
			accuracy = self.score.accuracy
			scoreMods = self.score.mods

			log.debug("[WIFIPIANO2] SCORE DATA: Stars: {stars}, OD: {od}, obj: {objects}, score: {score}, acc: {acc}, mods: {mods}".format(stars=stars, od=od, objects=objects, score=score, acc=accuracy, mods=scoreMods))

			# ---------- STRAIN PP
			# Scale score to mods multiplier
			scoreMultiplier = 1.0

			# Doubles score if EZ/HT
			if scoreMods & mods.EASY != 0:
				scoreMultiplier *= 0.50
			#if scoreMods & mods.HALFTIME != 0:
			#	scoreMultiplier *= 0.50

			# Calculate strain PP
			if scoreMultiplier <= 0:
				strainPP = 0
			else:
				score *= int(1.0 / scoreMultiplier)
				strainPP = pow(5.0 * max(1.0, stars / 0.0825) - 4.0, 3.0) / 110000.0
				strainPP *= 1 + 0.1 * min(1.0, float(objects) / 1500.0)
				if score <= 500000:
					strainPP *= (float(score) / 500000.0) * 0.1
				elif score <= 600000:
					strainPP *= 0.1 + float(score - 500000) / 100000.0 * 0.2
				elif score <= 700000:
					strainPP *= 0.3 + float(score - 600000) / 100000.0 * 0.35
				elif score <= 800000:
					strainPP *= 0.65 + float(score - 700000) / 100000.0 * 0.20
				elif score <= 900000:
					strainPP *= 0.85 + float(score - 800000) / 100000.0 * 0.1
				else:
					strainPP *= 0.95 + float(score - 900000) / 100000.0 * 0.05

			# ---------- ACC PP
			# Makes sure OD is in range 0-10. If this is done elsewhere, remove this.
			scrubbedOD = min(10.0, max(0, 10.0 - od))

			# Old formula but done backwards.
			hitWindow300 = (34 + 3 * scrubbedOD)

			# Increases hitWindow if EZ is on
			if scoreMods & mods.EASY != 0:
				hitWindow300 *= 1.4

			# Fiddles with DT and HT to make them match hitWindow300's ingame.
			if scoreMods & mods.DOUBLETIME != 0:
				hitWindow300 *= 1.5
			elif scoreMods & mods.HALFTIME != 0:
				hitWindow300 *= 0.75

			# makes hit window match what it is ingame.
			hitWindow300 = int(hitWindow300) + 0.5
			if scoreMods & mods.DOUBLETIME != 0:
				hitWindow300 /= 1.5
			elif scoreMods & mods.HALFTIME != 0:
				hitWindow300 /= 0.75

			# Calculate accuracy PP
			accPP = pow((150.0 / hitWindow300) * pow(accuracy, 16), 1.8) * 2.5
			accPP *= min(1.15, pow(float(objects) / 1500.0, 0.3))

			# ---------- TOTAL PP
			multiplier = 1.1
			if scoreMods & mods.NOFAIL != 0:
				multiplier *= 0.90
			if scoreMods & mods.SPUNOUT != 0:
				multiplier *= 0.95
			if scoreMods & mods.EASY != 0:
				multiplier *= 0.50
			pp = pow(pow(strainPP, 1.1) + pow(accPP, 1.1), 1.0 / 1.1) * multiplier
			log.debug("[WIFIPIANO2] Calculated PP: {}".format(pp))

			self.pp = pp
		except exceptions.invalidBeatmapException:
			log.warning("Invalid beatmap {}".format(self.beatmap.beatmapID))
			self.pp = 0
		finally:
			return self.pp