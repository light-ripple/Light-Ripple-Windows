import tornado.gen
import tornado.web

from common.web import requestsManager


class handler(requestsManager.asyncRequestHandler):
	@tornado.web.asynchronous
	@tornado.gen.engine
	def asyncGet(self):
		print("404: {}".format(self.request.uri))
		self.write("""
				<html>
					<head>
						<style>
							@import url(https://fonts.googleapis.com/css?family=Source+Sans+Pro:400,700,600,400italic,600italic,700italic,900,900italic);
							@import url(https://fonts.googleapis.com/css?family=Raleway:400,700);
							html, body {
								height: 90%;
								background-image: url(http://y.zxq.co/xtffuu.png);
							}
							.main {
								height: 100%;
								width: 100%;
								display: table;
							}
							.wrapper {
								display: table-cell;
								height: 90%;
								vertical-align: middle;
							}
							body {
								font-family: Source Sans Pro;
								text-align: center;
							}
							h1, h2, h3, h4, h5, h6 {
								font-family: Raleway;
							}
						</style>
					</head>
					<body>
						<div class = "main">
							<div class = "wrapper">
								<a href="https://ripple.moe"><img src="http://y.zxq.co/ufaibw.png"></a>
								<h3>Howdy, you're still connected to Ripple!</h3>
								You can't access osu!'s website if the Server Switcher is On.<br>
								Please open the <b>Server Switcher</b> and click <b>On/Off</b> to switch server, then refresh this page.
								<h4>If you still can't access osu! website even if the switcher is Off, <a href="http://www.refreshyourcache.com/" target="_blank">clean your browser cache</a>.</h4>
							</div>
						</div>
					</body>
				</html>
				""")
