if __name__ != "handlers.mods":
	from secret.achievements import common
	from objects import glob
	from common.constants import mods
else:
	import common

VERSION = 4
ORDER = 4

# Loads the achievement length on load
LENGTH = 0

ACHIEVEMENT_BASE = {
	"name": "{name}",
	"description": "{description}",
	"icon": "all-intro-{mod}"
}

ACHIEVEMENT_KEYS = {
	"name": [
		"Finality",
		"Perfectionist",
		"Rock Around The Clock",
		"Time And A Half",
		"Sweet Rave Party",
		"Blindsight",
		"Are You Afraid Of The Dark?",
		"Dial It Right Back",
		"Risk Averse",
		"Slowboat",
		"Burned Out"
	],
	"description": [
		"High stakes, no regrets.",
		"Accept nothing but the best.",
		"You can't stop the rock.",
		"Having a right ol' time. One and a half of them, almost.",
		"Founded in the fine tradition of changing things that were just fine as they were.",
		"I can see just perfectly.",
		"Harder than it looks, probably because it's hard to look.",
		"Sometimes you just want to take it easy.",
		"Safety nets are fun!",
		"You got there. Eventually.",
		"One cannot always spin to win."
	],
	"mod": [
		"suddendeath",
		"perfect",
		"hardrock",
		"doubletime",
		"nightcore",
		"hidden",
		"flashlight",
		"easy",
		"nofail",
		"halftime",
		"spunout"
	]
}

# For every iteration every index is increased by 1
ACHIEVEMENT_STRUCT = {
	"name": 1,
	"description": 1,
	"mod": 1
}

ACHIEVEMENTS = []

def load():
	global ACHIEVEMENTS, LENGTH
	ACHIEVEMENTS, LENGTH = common.load_achievement_data(ACHIEVEMENT_BASE, ACHIEVEMENT_KEYS, ACHIEVEMENT_STRUCT)

def handle(mode, score, beatmap, user_data):
	return check(score.mods)

def check(m):
	achievement_ids = []

	# Yes I am braindead atm and dont want to think about it...
	if m & mods.SUDDENDEATH > 0:
		achievement_ids += [0]
	if m & mods.PERFECT > 0:
		achievement_ids += [1]
	if m & mods.HARDROCK > 0:
		achievement_ids += [2]
	if m & mods.DOUBLETIME > 0:
		achievement_ids += [3]
	if m & mods.NIGHTCORE > 0:
		achievement_ids += [4]
	if m & mods.HIDDEN > 0:
		achievement_ids += [5]
	if m & mods.FLASHLIGHT > 0:
		achievement_ids += [6]
	if m & mods.EASY > 0:
		achievement_ids += [7]
	if m & mods.NOFAIL > 0:
		achievement_ids += [8]
	if m & mods.HALFTIME > 0:
		achievement_ids += [9]
	if m & mods.SPUNOUT > 0:
		achievement_ids += [10]

	return achievement_ids

def update(userID):
	achievement_ids = []

	entries = glob.db.fetchAll("SELECT mods FROM scores WHERE userid = %s GROUP BY mods", [userID])
	for entry in entries:
		achievement_ids += check(entry["mods"])

	return achievement_ids
