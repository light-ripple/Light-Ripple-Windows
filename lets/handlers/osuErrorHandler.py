import os

import tornado.gen
import tornado.web
import json

from common.web import requestsManager
from common.log import logUtils as log
from constants import exceptions
from common import generalUtils

MODULE_NAME = "clienterror"
class handler(requestsManager.asyncRequestHandler):
	@tornado.web.asynchronous
	@tornado.gen.engine
	def asyncPost(self): # Nothing here is useful for anyone -.-
		self.write("")
