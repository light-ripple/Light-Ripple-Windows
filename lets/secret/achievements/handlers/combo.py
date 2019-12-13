if __name__ != "handlers.combo":
	from secret.achievements import common
	from objects import glob
else:
	import common

VERSION = 1
ORDER = 1

# Loads the achievement length on load
LENGTH = 0

ACHIEVEMENT_BASE = {
	"name": "{index} Combo (osu!{mode})",
	"description": "{index} big ones! You're moving up in the world!",
	"icon": "osu-combo-{index}"
}

ACHIEVEMENT_KEYS = {
	"index": [500, 750, 1000, 2000],
	"mode": ["std", "taiko", "ctb", "mania"]
}

# For every iteration index gets increased and it will loop around when it hits its array length, while mode gets increased every 4 itterations
ACHIEVEMENT_STRUCT = {
	"index": 1,
	"mode": 4
}

ACHIEVEMENTS = []

def load():
	global ACHIEVEMENTS, LENGTH
	ACHIEVEMENTS, LENGTH = common.load_achievement_data(ACHIEVEMENT_BASE, ACHIEVEMENT_KEYS, ACHIEVEMENT_STRUCT)

def handle(mode, score, beatmap, user_data):
	return check(mode, score.maxCombo)

def check(mode, max_combo):
	achievement_ids = []
	indexies = [x for x in ACHIEVEMENT_KEYS["index"] if x <= max_combo]

	for index in range(len(indexies)):
		achievement_ids.append(index + mode * 4)

	return achievement_ids

def update(userID):
	achievement_ids = []

	entries = glob.db.fetchAll("SELECT MAX(max_combo) AS max_combo, play_mode FROM scores WHERE userid = %s AND completed >= 2 GROUP BY play_mode", [userID])
	for entry in entries:
		achievement_ids += check(entry["play_mode"], entry["max_combo"])

	return achievement_ids
