import requests
import json

from constants import exceptions
from objects import glob

def levbodRequest(handler, params=None):
	if params is None:
		params = {}
	result = requests.get("{}/{}".format(glob.conf.config["levbod"]["url"], handler), params=params)

	try:
		data = json.loads(result.text)
	except (json.JSONDecodeError, ValueError, requests.RequestException, KeyError, exceptions.noAPIDataError):
		return None

	if result.status_code != 200 or "data" not in data:
		return None

	return data["data"]

def getListing(rankedStatus=4, page=0, gameMode=-1, query=""):
	return levbodRequest("listing", {
		"mode": gameMode,
		"status": rankedStatus,
		"query": query,
		"page": page,
	})

def getBeatmapSet(id):
	return levbodRequest("beatmapset", {
		"id": id
	})

def getBeatmap(id):
	return levbodRequest("beatmap", {
		"id": id
	})

def levbodToDirect(data):
	s = "{beatmapset_id}.osz|{artist}|{title}|{creator}|{ranked_status}|10.00|0|{beatmapset_id}|".format(**data)
	if len(data["beatmaps"]) > 0:
		s += "{}|0|0|0||".format(data["beatmaps"][0]["beatmap_id"])
		for i in data["beatmaps"]:
			s += "{difficulty_name}@{game_mode},".format(**i)
	s = s.strip(",")
	s += "|"
	return s

def levbodToDirectNp(data):
	return "{beatmapset_id}.osz|{artist}|{title}|{creator}|{ranked_status}|10.00|0|{beatmapset_id}|{beatmapset_id}|0|0|0|".format(**data)