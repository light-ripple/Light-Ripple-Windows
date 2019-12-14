import os
import sys
import traceback
import json

import tornado.gen
import tornado.web
from raven.contrib.tornado import SentryMixin

from common.log import logUtils as log
from common.web import requestsManager
from constants import dataTypes
from constants import exceptions
from helpers import binaryHelper
from common import generalUtils
from objects import glob
from common.sentry import sentry

MODULE_NAME = "get_full_errorlog"
class handler(requestsManager.asyncRequestHandler):
	"""
	Handler for /errorlogs/
	"""
	@tornado.web.asynchronous
	@tornado.gen.engine
	@sentry.captureTornado
	def asyncGet(self, errorID):
		try:
			# Find the file
			fileName = ".data/clienterrors/{}.json".format(errorID)
			if not os.path.isfile(fileName):
				raise exceptions.fileNotFoundException(MODULE_NAME, fileName)
			
			# Read and json loads error file
			with open(fileName, "rb") as f:
				data = json.loads(f.read())
			
			# Write to page
			self.write(json.dumps(data, indent=2, sort_keys=True))
			self.add_header("Content-type", "application/json")
		except exceptions.fileNotFoundException:
			self.write("Errorlog not found")