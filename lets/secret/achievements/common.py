import math
if __name__ != "common":
	from objects import glob
	import time
	import json
	from common.ripple import userUtils

def load_achievement_data(ACHIEVEMENT_BASE, ACHIEVEMENT_KEYS, ACHIEVEMENT_STRUCT):
	LENGTH = 0
	ACHIEVEMENTS = []

	for struct in ACHIEVEMENT_STRUCT:
		LENGTH = max(LENGTH, len(ACHIEVEMENT_KEYS[struct]) * ACHIEVEMENT_STRUCT[struct])
	
	entry = {x:0 for x in ACHIEVEMENT_STRUCT}
	for i in range(LENGTH):
		for struct in ACHIEVEMENT_STRUCT:
			entry[struct] = math.floor(i / ACHIEVEMENT_STRUCT[struct]) % len(ACHIEVEMENT_KEYS[struct])
		format_data = {x:ACHIEVEMENT_KEYS[x][entry[x]] for x in ACHIEVEMENT_KEYS}
		ACHIEVEMENTS.append({x: ACHIEVEMENT_BASE[x].format_map(format_data) for x in ACHIEVEMENT_BASE})
	
	return ACHIEVEMENTS, LENGTH

def get_usercache(userID):
	user_cache = glob.redis.get("lets:user_achievement_cache:{}".format(userID))
	if user_cache is None:
		user_cache = {}
	else:
		user_cache = json.loads(user_cache.decode("utf-8"))

	if "version" not in user_cache:
		# Load from sql database
		user_cache["version"] = userUtils.getAchievementsVersion(userID)
		db_achievements = [x["achievement_id"] for x in glob.db.fetchAll("SELECT achievement_id FROM users_achievements WHERE user_id = %s", [userID])]
		if "achievements" in user_cache:
			user_cache["achievements"] += db_achievements
		else:
			user_cache["achievements"] = db_achievements
		# Remove duplicates after merge
		user_cache["achievements"] = list(set(user_cache["achievements"]))

	return user_cache

def add_pending_achievement(userID, achievementID):
	user_cache = get_usercache(userID)
	if len([x for x in user_cache["achievements"] if x in [achievementID, -achievementID]]) > 0:
		print("Tried to add achievement:{} to user:{}, but failed due to duplicate entry.".format(achievementID, userID))
		return

	user_cache["achievements"].append(-achievementID)

	# Remove duplicates after merge
	user_cache["achievements"] = list(set(user_cache["achievements"]))

	glob.redis.set("lets:user_achievement_cache:{}".format(userID), json.dumps(user_cache), 1800)

	userUtils.unlockAchievement(userID, achievementID)