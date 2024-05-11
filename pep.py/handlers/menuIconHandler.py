import tornado.web
import tornado.gen
import json

from common.web import requestsManager
from common.log import logUtils as log
from objects import banchoConfig as bc

class handler(requestsManager.asyncRequestHandler):
	#2024년에 https://assets.ppy.sh/menu-content.json 로 바뀜
	#Changed to https://assets.ppy.sh/menu-content.json in 2024
	@tornado.web.asynchronous
	@tornado.gen.engine
	def asyncGet(self):
		try:
			if self.request.host.startswith("assets"):
				imageURL = bc.banchoConfig.config["menuIcon"]
				if imageURL:
					url = imageURL.split("|")[1]
					imageURL = imageURL.split("|")[0]
					begins = expires = None
				else:
					imageURL = url = begins = expires = None

				data = {
					"images": [
						{
							"image": imageURL,
							"url": url,
							"IsCurrent": True,
							"begins": begins,
							"expires": expires
						}
					]
				}

				self.set_status(200)
				self.set_header('Content-Type', "application/json")
				self.write(json.dumps(data, indent=2))
			else:
				self.set_status(403)
				self.write('Host Is Not start "assets."')
		except Exception as e:
			log.error(f"menuIconHandler ERROR | {e}")
			self.set_status(500)