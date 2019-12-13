import json
import re
import traceback
from helpers import aeshelper
from objects import glob
from common.ripple import userUtils

from . import flavours
from . import ice_coffee
from . import police

#Cornflakes is nice when 90% is sugar
sugar = {
	"hash": [],
	"path": [],
	"file": [],
	"title": []
}

initialized_eggs = False

#Eggs
def init_eggs():
	eggs = glob.db.fetchAll("SELECT * FROM eggs", [])
	if eggs is not None:
		for egg in eggs:
			if egg["type"] not in ["hash", "path", "file", "title"]:
				continue
			sugar[egg["type"]].append(egg)

	compile_regex()

	initialized_eggs = True

def compile_regex():
	#Cache regex searches
	for carbohydrates in sugar:
		for speed in sugar[carbohydrates]:
			if speed["is_regex"]:
				speed["regex"] = re.compile(speed["value"])

#Since this is still being worked on everything is in a try catch
def bake(submit, score):
	"""
	We have deprecated the process list scanning.
	There is no config to re-enable this.
	If you know what you are doing however you know how to re-enable this feature.
	"""
	return
	try:
		if not initialized_eggs:
			init_eggs()
		
		if not score.passed:
			return
		
		detected = []
		flags = 0

		if "osuver" in submit.request.arguments:
			aeskey = "osu!-scoreburgr---------{}".format(submit.get_argument("osuver"))
		else:
			aeskey = "h89f2-890h2h89b34g-h80g134n90133"
		iv = submit.get_argument("iv")

		score_data = aeshelper.decryptRinjdael(aeskey, iv, submit.get_argument("score"), True).split(":")
		username = score_data[1].strip()

		user_id = userUtils.getID(username)
		restricted = userUtils.isRestricted(user_id)

		if restricted == True or user_id == 0: #We dont care about this since this person is already taken care off
			return

		flags = score_data[17].count(' ')

		try:
			pl = aeshelper.decryptRinjdael(aeskey, iv, submit.get_argument("pl"), True).split("\r\n")
		except:
			police.call("Unable to decrypt process list from USERNAME()", user_id=user_id)
			detected.append({
				"tag":"Unable to decrypt process list (Hacked)",
				"ban": False
				})
			eat(score, {}, detected, flags)
			return

		pl = sell(pl)

		#I dont really like chocolate that much >.<
		for p in pl:
			for t in sugar.keys():
				if p[t] is None:
					continue

				for speed in sugar[t]:
					if speed in detected:
						continue

					if speed["is_regex"]:
						if "regex" not in speed: #Some weird bug where it unsets itself
							speed["regex"] = re.compile(speed["value"])

						if speed["regex"].search(p[t]) is not None:
								detected.append(speed)
					else:
						if speed["value"] == p[t]:
							detected.append(speed)

		eat(score, pl, detected, flags)
	except:
		police.call(traceback.format_exc(), discord_m=True)
		police.call("Oh no! The cake is on fire! Abort!")

def sell(processes):
	formatted_pl = []
	for p in processes: #Formats the process list
		try:
			t = p.split(" | ", 1)
			try:
				d = t[0].split(" ", 1)
				file_hash = d[0]
				file_path = d[1]
			except:
				file_hash = None
				file_path = None

			h = t[1].split(" (", 1)
			file_name = h[0]

			file_title = None
			if len(h[1]) > 1:
				file_title = h[1][:-1]

			formatted_pl.append({"hash":file_hash, "path":file_path,
								 "file":file_name, "title":file_title})
		except:
			continue

	return formatted_pl

def eat(score, processes, detected, flags):
	if flavours.config is None:
		police.cache_config()

	do_restrict = False
	for toppings in detected:
		if toppings["ban"]:
			do_restrict = True

	tag_list = [x["tag"] for x in detected]

	hax_flags = flags & ~ice_coffee.IGNORE_HAX_FLAGS
	beatmap_id = get_beatmap_id(score.fileMd5)["beatmap_id"]
	
	username = userUtils.getUsername(score.playerUserID)

	fields = [
		{
			"name": "BeatmapID: {}".format(beatmap_id),
			"value": "[Download Beatmap](http://{}/b/{})".format(flavours.config["urls"]["main_domain"], beatmap_id),
			"inline": True
		},
		{
			"name": "ScoreID: {}".format(score.scoreID),
			"value": "[Download Replay](http://{}/web/replays/{})".format(flavours.config["urls"]["main_domain"], score.scoreID),
			"inline": True
		}
	]

	if len(detected) > 0:
		reason = " & ".join(tag_list)
		if len(reason) > 86:
			reason = "reasons..."

		extra_data = ""
		if hax_flags != 0:
			extra_data = "\nHad bad flags: ({}) -> ({})".format(flags, make_flags_string(flags))

		if do_restrict:
			userUtils.restrict(score.playerUserID)
			userUtils.appendNotes(score.playerUserID, "Restricted due to {}".format(reason))
			police.call("{} was restricted for: {} {}".format(username, reason, extra_data), 
				discord_m=True,
				embed_args={
						"color": 0xd9534f,
						"title": "Bad cake detected",
						"title_url": "http://old.{}/index.php?p=129&sid={}".format(flavours.config["urls"]["main_domain"], score.scoreID),
						"desc": "Restricted for: {} {}".format(reason, extra_data),
						"author": username,
						"author_icon": "http://a.{}/{}".format(flavours.config["urls"]["main_domain"], score.playerUserID),
						"author_url": "http://{}/u/{}".format(flavours.config["urls"]["main_domain"], score.playerUserID),
						"thumbnail": flavours.config["images"]["bad_cake_ban"],
						"fields": fields
					}
				)
		else:
			userUtils.appendNotes(score.playerUserID, reason)
			police.call("{} submitted bad cake: {} {}".format(username, reason, extra_data), 
				discord_m=True,
				embed_args={
						"color": 0xf0ad4e,
						"title": "Bad cake detected",
						"title_url": "http://old.{}/index.php?p=129&sid={}".format(flavours.config["urls"]["main_domain"], score.scoreID),
						"desc": "Had bad cake: {} {}".format(reason, extra_data),
						"author": username,
						"author_icon": "http://a.{}/{}".format(flavours.config["urls"]["main_domain"], score.playerUserID),
						"author_url": "http://{}/u/{}".format(flavours.config["urls"]["main_domain"], score.playerUserID),
						"thumbnail": flavours.config["images"]["bad_cake"],
						"fields": fields
					}
				)
	elif hax_flags != 0:
		police.call("{} submitted bad flags: ({}) -> ({})".format(username, flags, make_flags_string(flags)),
			discord_m=True, 
			embed_args={
					"color": 0xf0ad4e,
					"title": "Bad flags detected",
					"title_url": "http://old.{}/index.php?p=129&sid={}".format(flavours.config["urls"]["main_domain"], score.scoreID),
					"desc": "({}) -> ({})".format(flags, make_flags_string(flags)),
					"author": username,
					"author_icon": "http://a.{}/{}".format(flavours.config["urls"]["main_domain"], score.playerUserID),
					"author_url": "http://{}/u/{}".format(flavours.config["urls"]["main_domain"], score.playerUserID),
					"thumbnail": flavours.config["images"]["bad_flag"],
					"fields": fields
				}
			)

	glob.db.execute("INSERT INTO cakes(id, userid, score_id, processes, detected, flags) VALUES (NULL,%s,%s,%s,%s,%s)", [score.playerUserID, score.scoreID, json.dumps(processes), json.dumps(tag_list), flags])

def make_flags_string(i):
	s = []
	flags = [e for e in ice_coffee.Flags]

	for flag in flags:
		if i & flag.value and i & ~ice_coffee.IGNORE_HAX_FLAGS:
			s.append(flag.name)
	
	return ", ".join(s)

def get_beatmap_id(hash):
	query = "SELECT beatmap_id,beatmapset_id FROM beatmaps WHERE beatmap_md5 = %s"
	return glob.db.fetch(query, [hash])