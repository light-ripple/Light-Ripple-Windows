if __name__ != "__main__":
	print("This is ment to be runned as a toolkit for generating achievement database")
	exit()

import common
from os.path import dirname, basename, isfile
import glob
import importlib
import time
import math

SQL_STRING = """
CREATE TABLE IF NOT EXISTS `achievements` (
  `id` int(11) NOT NULL,
  `name` varchar(64) NOT NULL,
  `description` varchar(128) NOT NULL,
  `icon` varchar(32) NOT NULL,
  `version` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
INSERT INTO achievements (id, name, description, icon, version) VALUES
"""

module_list = glob.glob("handlers/*.py")
module_list = [basename(f)[:-3] for f in module_list if isfile(f) and not f.endswith("__init__.py")]
#							^ cat face

modules = []
for module in module_list:
	modules.append(importlib.import_module("handlers.{}".format(module)))

modules = sorted(modules, key=lambda k: k.ORDER)

SQL_INSERTS = []
index = 1
for module in modules:
	module.load()
	for achievement in module.ACHIEVEMENTS:
		SQL_INSERTS.append("({}, '{}', '{}', '{}', {})".format(
			index,
			achievement["name"].replace('"', '\\"').replace("'", "\\'"),
			achievement["description"].replace('"', '\\"').replace("'", "\\'"),
			achievement["icon"].replace('"', '\\"').replace("'", "\\'"),
			module.VERSION
			))
		index += 1

SQL_STRING += ",\n".join(SQL_INSERTS) + ";"

FILENAME = "achievements-{}.sql".format(math.floor(time.time()))
with open(FILENAME, "w") as f:
	f.write(SQL_STRING)

print("Saved sql export into {}".format(FILENAME))
print("Import this table into your database.")
print("""NOTE: Avoid changing the ORDER variable inside the handlers at all cost as this will result in
new achievement sql data not matching data of already achieved achievements by users.""")
print("If you know what you are doing you know how to fix this if you still choose to ignore this warning")