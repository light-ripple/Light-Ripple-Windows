from common.constants import bcolors
from objects import glob
from random import choice
colors = ["\033[95m", "\033[94m", "\033[92m", "\033[93m", "\033[91m"]
color = (choice((colors)))

def printServerStartHeader(asciiArt=True):

	if asciiArt:
		print("{}           _                 __".format(color))
		print("          (_)              /  /")
		print("   ______ __ ____   ____  /  /____")
		print("  /  ___/  /  _  \\/  _  \\/  /  _  \\")
		print(" /  /  /  /  /_) /  /_) /  /  ____/")
		print("/__/  /__/  .___/  .___/__/ \\_____/")
		print("        /  /   /  /")
		print("       /__/   /__/\r\n")
		print("                          .. o  .")
		print("                         o.o o . o")
		print("                        oo...")
		print("                    __[]__")
		print("             ______/o_o_o_|__       u wot m8? ")
		print("             \\\"\"\"\"\"\"\"\"\"\"\"\"\"\"/")
		print("              \\ . ..  .. . /")
		print("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	printColored("> Welcome to pep.py v2 | The osu!bancho server for Windows", bcolors.GREEN)
	printColored("> Made by the Ripple team | Rewrite by Uniminin", bcolors.GREEN)
	printColored("> {}https://zxq.co/ripple/pep.py".format(bcolors.UNDERLINE), bcolors.GREEN)
	printColored("> Press CTRL+C to exit | Kill Console In Windows to exit.\n", bcolors.GREEN)

def printNoNl(string):
	"""
	Print a string without \n at the end

	:param string: string to print
	:return:
	"""
	print(string, end="")

def printColored(string, color):
	"""
	Print a colored string

	:param string: string to print
	:param color: ANSI color code
	:return:
	"""
	print("{}{}{}".format(color, string, bcolors.ENDC))

def printError():
	"""
	Print a red "Error"

	:return:
	"""
	printColored("Error", bcolors.RED)

def printDone():
	"""
	Print a green "Done"

	:return:
	"""
	printColored("Done", bcolors.GREEN)

def printWarning():
	"""
	Print a yellow "Warning"

	:return:
	"""
	printColored("Warning", bcolors.YELLOW)
