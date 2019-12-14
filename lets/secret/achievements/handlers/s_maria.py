if __name__ != "handlers.s_maria":
	from secret.achievements import common
	from objects import glob
else:
	import common

VERSION = 6
ORDER = 6

# Loads the achievement length on load
LENGTH = 0

ACHIEVEMENT_BASE = {
	"name": "{name}",
	"description": "{description}",
	"icon": "{icon}"
}

ACHIEVEMENT_KEYS = {
	"name": ["A meganekko approaches"],
	"description": ["Congratulations, you met Maria!"],
    "icon": ["mania-secret-meganekko"]
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
	return check(mode, score.maxCombo)

def check(mode, max_combo):
	if mode is not 3 or max_combo < 100:
		return []
	return [0]

def update(userID):
	found = glob.db.fetch("SELECT id FROM scores WHERE userid = %s AND play_mode = 3 AND completed >= 2 AND max_combo >= 100 LIMIT 1;", [userID])
	if found is None:
		return []
	return [0]
