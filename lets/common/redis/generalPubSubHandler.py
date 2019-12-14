import json

def shape(d):
	"""
	Returns a shape of a dictionary.
	Used to check if two dictionaries have the same structure

	:param d: dictionary
	:return: `d`'s shape
	"""
	if isinstance(d, dict):
		return {k: shape(d[k]) for k in d}
	else:
		return None

class wrongStructureError(Exception):
	pass

class generalPubSubHandler:
	def __init__(self):
		self.structure = {}
		self.type = "json"
		self.strict = True

	def parseData(self, data):
		"""
		Parse received data

		:param data: received data, as bytes array
		:return: parsed data or None if it's invalid
		"""
		if self.type == "json":
			# Parse json
			if type(data) == int:
				return None
			data = json.loads(data.decode("utf-8"))
			if shape(data) != shape(self.structure) and self.strict:
				raise wrongStructureError()
		elif self.type == "int":
			# Parse int
			data = int(data.decode("utf-8"))
		return data