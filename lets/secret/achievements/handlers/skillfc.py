if __name__ != "handlers.skillfc":
	import math
	from secret.achievements import common
	from common.ripple import scoreUtils
	from objects import glob, beatmap
else:
	import common

VERSION = 3
ORDER = 3

# Loads the achievement length on load
LENGTH = 0

ACHIEVEMENT_BASE = {
	"name": "{name}",
	"description": "{description}",
	"icon": "{mode}-skill-fc-{index}"
}

ACHIEVEMENT_KEYS = {
	"index": [1, 2, 3, 4, 5, 6, 7, 8],
	"mode": ["osu", "taiko", "fruits", "mania"],
	"name": [
		"Totality",
		"Keeping Time",
		"Sweet And Sour",
		"Keystruck",
		"Business As Usual",
		"To Your Own Beat",
		"Reaching The Core",
		"Keying In",
		"Building Steam",
		"Big Drums",
		"Clean Platter",
		"Hyperflow",
		"Moving Forward",
		"Adversity Overcome",
		"Between The Rain",
		"Breakthrough",
		"Paradigm Shift",
		"Demonslayer",
		"Addicted",
		"Everything Extra",
		"Anguish Quelled",
		"Rhythm's Call",
		"Quickening",
		"Level Breaker",
		"Never Give Up",
		"Time Everlasting",
		"Supersonic",
		"Step Up",
		"Aberration",
		"The Drummer's Throne",
		"Dashing Scarlet",
		"Behind The Veil"
	],
	"description": [
		"All the notes. Every single one.",
		"Two to go, please.",
		"Hey, this isn't so bad.",
		"Bet you feel good about that.",
		"Surprisingly difficult.",
		"Don't choke.",
		"Excellence is its own reward.",
		"They said it couldn't be done. They were wrong."
	]
}

# For every iteration name and mode is increased by 1 and loops around when they hit their length while index and description increases by 1 every 4 itterations
ACHIEVEMENT_STRUCT = {
	"name": 1,
	"mode": 1,
	"index": 4,
	"description": 4
}

ACHIEVEMENTS = []

def load():
	global ACHIEVEMENTS, LENGTH
	ACHIEVEMENTS, LENGTH = common.load_achievement_data(ACHIEVEMENT_BASE, ACHIEVEMENT_KEYS, ACHIEVEMENT_STRUCT)

def handle(mode, score, beatmap, user_data):
	if not score.fullCombo: # No need to check if the score were not a fullcombo
		return []
	return check(mode, beatmap)

def check(mode, beatmap):
	achievement_ids = []

	mode_str = scoreUtils.readableGameMode(mode)

	mode_2 = mode_str.replace("osu", "std")
	stars = getattr(beatmap, "stars" + mode_2.title())

	indexies = [x - 1 for x in ACHIEVEMENT_KEYS["index"] if x == math.floor(stars)]

	for index in indexies:
		achievement_ids.append(mode + index * 4)

	return achievement_ids

def update(userID):
	achievement_ids = []

	entries = glob.db.fetchAll("SELECT beatmap_md5, play_mode FROM scores WHERE full_combo = 1 AND completed >= 2 AND userid = %s GROUP BY beatmap_md5, play_mode", [userID])
	for entry in entries:
		current_beatmap = beatmap.beatmap()
		current_beatmap.setDataFromDB(entry["beatmap_md5"])
		
		achievement_ids += check(entry["play_mode"], current_beatmap)

	return achievement_ids
