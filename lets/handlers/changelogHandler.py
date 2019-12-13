from datetime import datetime
 
import tornado.gen
import tornado.web

from common.sentry import sentry
from common.web import requestsManager


class ChangelogDate:
    def __init__(self, timestamp):
        self.timestamp = timestamp

    def __str__(self):
        return datetime.fromtimestamp(self.timestamp).strftime("%b %d, %Y")


class ChangelogEntry:
    def __init__(self, timestamp, author, description):
        self.timestamp = ChangelogDate(int(timestamp))
        self.author = author.strip()
        self._description = description.strip()
        #self.repo = repo.strip()

    @property
    def description(self):
        return self._description.lstrip("*").strip().lstrip("+").strip().replace("🔺", "^").replace("🔼", "^")

    @property
    def symbol(self):
        return \
            "*" if self.description.startswith(("Fix", "*")) else \
            "+" if self.description.startswith(("Add", "+")) else \
            ""

    def __str__(self):
         return f"{self.symbol}\t{self.author}\t: {self.description}"


class handler(requestsManager.asyncRequestHandler):
    @tornado.web.asynchronous
    @tornado.gen.engine
    @sentry.captureTornado
    def asyncGet(self):
        output = ""

        try:
            lines = []
            with open("../changelog.txt", "r") as f:
                for i, l in enumerate(f):
                    if i >= 100:
                        break
                    lines.append(l)

            changelog_entries = []
            for line in lines:
                parts = line.split("|")
                if len(parts) != 4:
                    continue
                changelog_entries.append(ChangelogEntry(*parts[1:]))

            if not changelog_entries:
                return
            last_day = changelog_entries[0].timestamp.timestamp // 86400
            output = f"{str(changelog_entries[0].timestamp)}"
            for i, entry in enumerate(changelog_entries):
                this_day = entry.timestamp.timestamp // 86400
                if this_day != last_day:
                    last_day = this_day
                    output += f"\n{str(entry.timestamp)}"
                output += f"\n{str(entry)}"
        finally:
            self.write(output)
