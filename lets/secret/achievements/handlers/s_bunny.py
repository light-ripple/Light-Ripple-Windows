if __name__ != "handlers.s_bunny":
	from secret.achievements import common
	from objects import glob
else:
	import common

VERSION = 6
ORDER = 7

# Loads the achievement length on load
LENGTH = 0

ACHIEVEMENT_BASE = {
	"name": "{name}",
	"description": "{description}",
	"icon": "{icon}"
}

ACHIEVEMENT_KEYS = {
	"name": ["Don't let the bunny distract you!"],
	"description": ["The order was indeed, not a rabbit."],
	"icon": ["all-secret-bunny"]
}

# Might implement something to "automate" this if there is only one achievement for this handler (Make different structure)
ACHIEVEMENT_STRUCT = {
	"name": 1,
	"description": 1,
	"icon": 1
}

ACHIEVEMENTS = []

def load():
	global ACHIEVEMENTS, LENGTH
	ACHIEVEMENTS, LENGTH = common.load_achievement_data(ACHIEVEMENT_BASE, ACHIEVEMENT_KEYS, ACHIEVEMENT_STRUCT)

def handle(mode, score, beatmap, user_data):
	return check(mode, score, beatmap.beatmapSetID)

def check(mode, score, beatmapSetID):
	if mode is not 0 or beatmapSetID is not 184 or not score.fullCombo:
		return []
	return [0]

def update(userID):
	found = glob.db.fetch("SELECT id FROM scores WHERE userid = %s AND play_mode = 0 AND beatmap_md5 IN ('6fe6d6e6d2fae3da30140b640c7ec49a', '642727d50438686f6f35cc1aacbcb3cd') AND completed >= 2 AND full_combo = 1 LIMIT 1;", [userID])
	if found is None:
		return []
	return [0]
