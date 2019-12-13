var parser = require("gettext-parser")
var fs     = require("fs")

fs.readFile("data/js-locales/templates.pot", "utf-8", (err, data) => {
	if (err)
		throw err
	var tpl = parser.po.parse(data, "utf-8")
	var plurals = {}
	Object.keys(tpl.translations[""]).forEach(key => {
		let val = tpl.translations[""][key]
		
		if (typeof val === "undefined")
			return
		
		if (!val.msgid_plural)
			return
		
		plurals[val.msgid] = val.msgid_plural
	})

	fs.writeFile("static/key_plural.js", "var keyPlurals = " + JSON.stringify(plurals) + ";", err => {
		if (err) throw err;
	})
})
