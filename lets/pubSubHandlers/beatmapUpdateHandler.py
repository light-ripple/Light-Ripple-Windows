from common.redis import generalPubSubHandler
from helpers import osuapiHelper
from objects import beatmap

def updateSet(beatmapSetID):
	apiResponse = osuapiHelper.osuApiRequest("get_beatmaps", "s={}".format(beatmapSetID), False)
	if len(apiResponse) == 0:
		return
	for i in apiResponse:
		beatmap.beatmap(i["file_md5"], int(i["beatmapset_id"]), refresh=True)

class handler(generalPubSubHandler.generalPubSubHandler):
	def __init__(self):
		super().__init__()
		self.structure = {}
		self.strict = False

	def handle(self, data):
		data = super().parseData(data)
		if data is None:
			return
		if "id" in data:
			beatmapData = osuapiHelper.osuApiRequest("get_beatmaps", "b={}".format(data["id"]))
			if beatmapData is not None and "beatmapset_id" in beatmapData:
				updateSet(beatmapData["beatmapset_id"])
		elif "set_id" in data:
			updateSet(data["set_id"])