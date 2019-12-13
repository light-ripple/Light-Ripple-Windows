import os
import sys
import traceback
import string
import time
import hashlib

import tornado.gen
import tornado.web

from common.log import logUtils as log
from common.ripple import userUtils
from common.ripple import passwordUtils
from common.web import requestsManager
from common import generalUtils
from objects import glob
from common.sentry import sentry


MODULE_NAME = "inGameRegistrationHandler"
class handler(requestsManager.asyncRequestHandler):
	"""
	Handler for /users
		by Cherry and Aoba
	"""
	@tornado.web.asynchronous
	@tornado.gen.engine
	def asyncPost(self):
		try:
			if not requestsManager.checkArguments(self.request.arguments, ["user[username]", "user[user_email]", "user[password]", "check"]):
				return self.write("what are you doing here?")
			username = self.get_argument("user[username]")
			email = self.get_argument("user[user_email]")
			password = self.get_argument("user[password]") # Raw password
			accountCreated = self.get_argument("check")
			if accountCreated == "1":
				return self.write('{"form_error":{"user":{"check":["Account already created."]}}}') 
			emailCheck = glob.db.fetch("SELECT 1 FROM users WHERE email = %s", [email])
			usernameCheck = glob.db.fetch("SELECT 1 FROM users WHERE username = %s", [username])	
			if emailCheck != None:
				return self.write('{"form_error":{"user":{"user_email":["Email address already used."]}}}')
			if usernameCheck != None or username.lower() in ["peppy","rrtyui","cookiezi","azer","loctav","banchobot","happystick","doomsday","sharingan33","andrea","cptnxn","reimu-desu","hvick225","_index","my aim sucks","kynan","rafis","sayonara-bye","thelewa","wubwoofwolf","millhioref","tom94","tillerino","clsw","spectator","exgon","axarious","angelsim","recia","nara","emperorpenguin83","bikko","xilver","vettel","kuu01","_yu68","tasuke912","dusk","ttobas","velperk","jakads","jhlee0133","abcdullah","yuko-","entozer","hdhr","ekoro","snowwhite","osuplayer111","musty","nero","elysion","ztrot","koreapenguin","fort","asphyxia","niko","shigetora"]:
				return self.write('{"form_error":{"user":{"username":["Username already used or it is forbidden."]}}}')	
			if len(password) < 8 or len(password) > 32:
				return self.write('{"form_error":{"user":{"password":["Password too short or long! (Password length must be more than 8 and less than 32)"]}}}') 
			if "_" in username and " " in username:
				self.write('{"form_error":{"user":{"username":["An username can not contain both underscores and spaces."]}}}')
			userID = int(glob.db.execute("INSERT INTO users(username, username_safe, password_md5, salt, email, register_datetime, privileges, password_version) VALUES (%s,       %s,            %s,            '',  %s,     %s,                 1048576,          2)", [username, userUtils.safeUsername(username), passwordUtils.genBcrypt(hashlib.md5(password.encode('utf-8')).hexdigest()), email, int(time.time())]))
			glob.db.execute("INSERT INTO users_stats(id, username, user_color, user_style, ranked_score_std, playcount_std, total_score_std, ranked_score_taiko, playcount_taiko, total_score_taiko, ranked_score_ctb, playcount_ctb, total_score_ctb, ranked_score_mania, playcount_mania, total_score_mania) VALUES (%s, %s, 'black', '', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)", [userID, username])
			glob.db.execute("INSERT INTO rx_stats(id, username, user_color, user_style, ranked_score_std, playcount_std, total_score_std, ranked_score_taiko, playcount_taiko, total_score_taiko, ranked_score_ctb, playcount_ctb, total_score_ctb, ranked_score_mania, playcount_mania, total_score_mania) VALUES (%s, %s, 'black', '', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)", [userID, username])
			log.info("{} created their account using ingame registration.".format(username))
		except Exception as e:
			log.error(e)