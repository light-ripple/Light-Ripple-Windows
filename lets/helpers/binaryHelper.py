"""That's basically packetHelper.py from pep.py, with some changes to make it work with replay files."""

from constants import dataTypes
import struct

def uleb128Encode(num):
	arr = bytearray()
	length = 0
	if num == 0:
		return bytearray(b"\x00")
	while num > 0:
		arr.append(num & 127)
		num >>= 7
		if num != 0:
			arr[length] |= 128
		length+=1
	return arr

def packData(__data, __dataType):
	data = bytes()
	pack = True
	packType = "<B"
	if __dataType == dataTypes.bbytes:
		pack = False
		data = __data
	elif __dataType == dataTypes.string:
		pack = False
		if len(__data) == 0:
			data += b"\x00"
		else:
			data += b"\x0B"
			data += uleb128Encode(len(__data))
			data += str.encode(__data, "latin_1")
	elif __dataType == dataTypes.uInt16:
		packType = "<H"
	elif __dataType == dataTypes.sInt16:
		packType = "<h"
	elif __dataType == dataTypes.uInt32:
		packType = "<L"
	elif __dataType == dataTypes.sInt32:
		packType = "<l"
	elif __dataType == dataTypes.uInt64:
		packType = "<Q"
	elif __dataType == dataTypes.sInt64:
		packType = "<q"
	elif __dataType == dataTypes.string:
		packType = "<s"
	elif __dataType == dataTypes.ffloat:
		packType = "<f"
	elif __dataType == dataTypes.rawReplay:
		pack = False
		data += packData(len(__data), dataTypes.uInt32)
		data += __data
	if pack:
		data += struct.pack(packType, __data)
	return data


def binaryWrite(structure = None):
	if structure is None:
		structure = []
	packetData = bytes()
	for i in structure:
		packetData += packData(i[0], i[1])
	return packetData
