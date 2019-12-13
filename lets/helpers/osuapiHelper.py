import json
from urllib.parse import quote

import requests

from common.log import logUtils as log
from common import generalUtils
from objects import glob
from constants import exceptions


def osuApiRequest(request, params, getFirst=True):
	"""
	Send a request to osu!api.

	request -- request type, string (es: get_beatmaps)
	params -- GET parameters, without api key or trailing ?/& (es: h=a5b99395a42bd55bc5eb1d2411cbdf8b&limit=10)
	return -- dictionary with json response if success, None if failed or empty response.
	"""
	# Make sure osuapi is enabled
	if not generalUtils.stringToBool(glob.conf.config["osuapi"]["enable"]):
		log.warning("osu!api is disabled")
		return None

	# Api request
	resp = None
	try:
		finalURL = "{}/api/{}?k={}&{}".format(glob.conf.config["osuapi"]["apiurl"], request, glob.conf.config["osuapi"]["apikey"], params)
		log.debug(finalURL)
		resp = requests.get(finalURL, timeout=5).text
		data = json.loads(resp)
		if getFirst:
			if len(data) >= 1:
				resp = data[0]
			else:
				resp = None
		else:
			resp = data
	finally:
		glob.dog.increment(glob.DATADOG_PREFIX+".osu_api.requests")
		log.debug(str(resp).encode("utf-8"))
		return resp

def getOsuFileFromName(fileName):
	"""
	Send a request to osu! servers to download a .osu file from file name
	Used to update beatmaps

	fileName -- .osu file name to download
	return -- .osu file content if success, None if failed
	"""
	# Make sure osuapi is enabled
	if not generalUtils.stringToBool(glob.conf.config["osuapi"]["enable"]):
		log.warning("osuapi is disabled")
		return None

	response = None
	try:
		URL = "{}/web/maps/{}".format(glob.conf.config["osuapi"]["apiurl"], quote(fileName))
		req = requests.get(URL, timeout=20)
		req.encoding = "utf-8"
		response = req.content
	finally:
		glob.dog.increment(glob.DATADOG_PREFIX+".osu_api.osu_file_requests")
		return response

def getOsuFileFromID(beatmapID):
	"""
	Send a request to osu! servers to download a .osu file from beatmap ID
	Used to get .osu files for oppai

	beatmapID -- ID of beatmap (not beatmapset) to download
	return -- .osu file content if success, None if failed
	"""
	# Make sure osuapi is enabled
	if not generalUtils.stringToBool(glob.conf.config["osuapi"]["enable"]):
		log.warning("osuapi is disabled")
		return None

	response = None
	try:
		URL = "{}/osu/{}".format(glob.conf.config["osuapi"]["apiurl"], beatmapID)
		response = requests.get(URL, timeout=20).content
	finally:
		glob.dog.increment(glob.DATADOG_PREFIX+".osu_api.osu_file_requests")
		return response