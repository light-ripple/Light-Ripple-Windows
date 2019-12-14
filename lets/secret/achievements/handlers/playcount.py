if __name__ != "handlers.playcount":
	from secret.achievements import common
	from objects import glob
else:
	import common

VERSION = 5
ORDER = 5

# Loads the achievement length on load
LENGTH = 0

ACHIEVEMENT_BASE = {
	"name": "{index_formatted} Plays",
	"description": "{description}",
	"icon": "osu-plays-{index}"
}

ACHIEVEMENT_KEYS = {
	"index": [5000, 15000, 25000, 50000],
	"index_formatted": ["5,000", "15,000", "25,000", "50,000"],
	"description": [
		"There's a lot more where that came from.",
		"Must.. click.. circles..",
		"There's no going back.",
		"You're here forever."
	]
}

# For every iteration every index is increased by 1
ACHIEVEMENT_STRUCT = {
	"index": 1,
	"index_formatted": 1,
	"description": 1
}

ACHIEVEMENTS = []

def load():
	global ACHIEVEMENTS, LENGTH
	ACHIEVEMENTS, LENGTH = common.load_achievement_data(ACHIEVEMENT_BASE, ACHIEVEMENT_KEYS, ACHIEVEMENT_STRUCT)

def handle(mode, score, beatmap, user_data):
	if mode is not 0:
		return []
	return check(user_data["playcount"])

def check(playcount):
	achievement_ids = []
	indexies = [x for x in ACHIEVEMENT_KEYS["index"] if x <= playcount]

	for index in range(len(indexies)):
		achievement_ids.append(index)

	return achievement_ids

def update(userID):
	achievement_ids = []

	playcount = glob.db.fetch("SELECT playcount_std FROM users_stats WHERE id = %s", [userID])["playcount_std"]
	achievement_ids += check(playcount)

	return achievement_ids
