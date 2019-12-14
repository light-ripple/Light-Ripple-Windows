import sys
import traceback
from functools import wraps

from common.log import logUtils as log


def trackExceptions(moduleName=""):
	def _trackExceptions(func):
		def _decorator(request, *args, **kwargs):
			try:
				response = func(request, *args, **kwargs)
				return response
			except:
				log.error("Unknown error{}!\n```\n{}\n{}```".format(" in "+moduleName if moduleName != "" else "", sys.exc_info(), traceback.format_exc()), True)
		return wraps(func)(_decorator)
	return _trackExceptions
