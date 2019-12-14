import sys
import traceback
import json

import tornado.gen
import tornado.web
from raven.contrib.tornado import SentryMixin

from common.log import logUtils as log
from common.ripple import userUtils
from common.web import requestsManager
from constants import exceptions
from objects import glob
from common.sentry import sentry

MODULE_NAME = "chart_handler"
class handler(requestsManager.asyncRequestHandler):
	"""
	Handler for /web/osu-getcharts.php
	"""
    @tornado.web.asynchronous
	@tornado.gen.engine
	@sentry.captureTornado
	def asyncGet(self):
        try:
        # Argument check
			if not requestsManager.checkArguments(self.request.arguments, ["u", "h"]):
				raise exceptions.invalidArgumentsException(MODULE_NAME)
            username = self.get_argument("u")
			password = self.get_argument("h")

            localUserId = glob.db.fetch("SELECT * FROM users WHERE username = %s", username)
            if localUserId < 0:
                return
            
            charts = glob.db.fetchAll("SELECT * FROM charts WHERE active = 1 ORDER BY end_date DESC LIMIT 10")
            for counter in range(0, len(charts)):
                charts['beatmaps'] = glob.db.fetchAll("SELECT * FROM scores INNER JOIN beatmaps ON score.beatmap_md5 = beatmaps.beatmap_md5 WHERE chart_id = %s ORDER BY pp DESC")
            return json.dumps(charts)
        except exceptions.invalidArgumentsException:
			pass