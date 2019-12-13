from objects import glob
from common.ripple import userUtils, scoreUtils
from os.path import dirname, basename, isfile
import glob as _glob
import importlib
import json
from secret.achievements import common

def load_achievements():
	"""Load all the achievements from handler list into glob.achievementClasses,
	and sets glob.ACHIEVEMENTS_VERSION to the highest version number in our achievement list.
	"""

	modules = _glob.glob("secret/achievements/handlers/*.py")
	modules = [basename(f)[:-3] for f in modules if isfile(f) and not f.endswith("__init__.py")]
	#						^ cat face

	for module in modules:
		module = importlib.import_module("secret.achievements.handlers." + module)
		module.load()
		if module.ORDER in glob.achievementClasses:
			print("\n!!! FOUND OVERLAPPING ACHIEVEMENT ORDER FOR {}!!!".format(module.ORDER))
			print("Unable to load {} due to {} already loaded in slot {}\n".format(module.__name__, glob.achievementClasses[module.ORDER].__name__, module.ORDER))
			continue
		glob.achievementClasses[module.ORDER] = module
		glob.ACHIEVEMENTS_VERSION = max(glob.ACHIEVEMENTS_VERSION, module.VERSION)
	
	print("Loaded {} achievement classes!".format(len(glob.achievementClasses)), end=" ")

def unlock_achievements_update(userID, version):
	"""Scans the user for past achievements they should have unlocked
	
	Arguments:
		userID {int} -- User id of a player
		version {int} -- Last achivement version the player had
	
	Returns:
		Array -- List of achievements
	"""
	achievements = []

	# Scan all past achivement versions from the user's achivement version to the latest
	index = 1
	for handler in glob.achievementClasses.values():
		if handler.VERSION > version:
			achievements += [x + index for x in handler.update(userID)]
		index += handler.LENGTH

	# Update achivement version for user
	userUtils.updateAchievementsVersion(userID)

	return achievements

def unlock_achievements(score, beatmap, user_data):
	"""Return array of achievements the current play recived
	
	Arguments:
		score {Score} -- Score data recived from replay
		beatmap {Beatmap} -- Played beatmap
		user_data {dict} -- Info about the current player
	
	Returns:
		Array -- List of achievements for the current play
	"""
	# Check if we have a valid mods
	if not scoreUtils.isRankable(score.mods):
		return []

	achievements = []

	userID = userUtils.getID(score.playerName)
	user_cache = common.get_usercache(userID)

	# Get current gamemode and change value std to osu
	gamemode_index = score.gameMode

	# Check if user should run achivement recheck
	if user_cache["version"] < glob.ACHIEVEMENTS_VERSION:
		achievements += unlock_achievements_update(userID, user_cache["version"])

	# Check if gameplay should get new achivement
	index = 1
	for handler in glob.achievementClasses.values():
		achievements += [x + index for x in handler.handle(gamemode_index, score, beatmap, user_data)]
		index += handler.LENGTH
	
	# Add pending achievements that were added though redis or mysql
	achievements += [-x for x in user_cache["achievements"] if x < 0] # Negative achievements id's means its pending

	# Remove pending achievements from redis object since we added it to the post achievements
	user_cache["achievements"] = [x for x in user_cache["achievements"] if x > 0]

	# Remove duplicated achievements (incase of unlock_achievements_update adding stuff)
	achievements = list(set(achievements))

	# Remove already achived achievements from list
	achievements = [x for x in achievements if x not in user_cache["achievements"]]

	user_cache["achievements"] += achievements
	glob.redis.set("lets:user_achievement_cache:{}".format(userID), json.dumps(user_cache), 1800)

	for achievement in achievements:
		userUtils.unlockAchievement(userID, achievement)

	return achievements

def achievements_response(achievements):
	achievement_objects = []

	index = 1
	for handler in glob.achievementClasses.values():
		achievement_objects += [handler.ACHIEVEMENTS[x - index] for x in achievements if len(handler.ACHIEVEMENTS) > x - index and x - index >= 0]
		index += handler.LENGTH

	achievements_packed = []
	for achievement_object in achievement_objects:
		achievements_packed.append("+".join([achievement_object["icon"], achievement_object["name"], achievement_object["description"]]))

	return "/".join(achievements_packed)
