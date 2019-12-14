if __name__ != "handlers.skillpass":
	import math
	from secret.achievements import common
	from common.ripple import scoreUtils
	from objects import glob, beatmap
else:
	import common

VERSION = 2
ORDER = 2

# Loads the achievement length on load
LENGTH = 0

ACHIEVEMENT_BASE = {
	"name": "{name}",
	"description": "{description}",
	"icon": "{mode}-skill-pass-{index}"
}

ACHIEVEMENT_KEYS = {
	"index": [1, 2, 3, 4, 5, 6, 7, 8],
	"mode": ["osu", "taiko", "fruits", "mania"],
	"name": [
		"Rising Star",
		"My First Don",
		"A Slice Of Life",
		"First Steps",
		"Constellation Prize",
		"Katsu Katsu Katsu",
		"Dashing Ever Forward",
		"No Normal Player",
		"Building Confidence",
		"Not Even Trying",
		"Zesty Disposition",
		"Impulse Drive",
		"Insanity Approaches",
		"Face Your Demons",
		"Hyperdash ON!",
		"Hyperspeed",
		"These Clarion Skies",
		"The Demon Within",
		"It's Raining Fruit",
		"Ever Onwards",
		"Above and Beyond",
		"Drumbreaker",
		"Fruit Ninja",
		"Another Surpassed",
		"Supremacy",
		"The Godfather",
		"Dreamcatcher",
		"Extra Credit",
		"Absolution",
		"Rhythm Incarnate",
		"Lord of the Catch",
		"Maniac"
	],
	"description": [
		"Can't go forward without the first steps.",
		"Definitely not a consolation prize. Now things start getting hard!",
		"Oh, you've SO got this.",
		"You're not twitching, you're just ready.",
		"Everything seems so clear now.",
		"A cut above the rest.",
		"All marvel before your prowess.",
		"My god, you're full of stars!"
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

	entries = glob.db.fetchAll("SELECT beatmap_md5, play_mode FROM scores WHERE completed = 3 AND userid = %s", [userID])
	for entry in entries:
		current_beatmap = beatmap.beatmap()
		current_beatmap.setDataFromDB(entry["beatmap_md5"])

		achievement_ids += check(entry["play_mode"], current_beatmap)

	return achievement_ids
