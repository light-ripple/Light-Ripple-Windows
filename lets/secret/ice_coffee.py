from enum import Enum

class Flags(Enum):
    CLEAN = 0
    SPEED = 1 << 1
    INCORRECT_MOD = 1 << 2
    MULTIPLE_OSU_CLIENTS = 1 << 3
    CHECKSUM_FAIL = 1 << 4
    FLASHLIGHT_CHECKSUM_FAIL = 1 << 5
    OSU_CHECKSUM = 1 << 6 #We already know...
    MISSING_PL = 1 << 7 #We already know...
    FLASHLIGHT_IMAGE = 1 << 8
    SPINNER = 1 << 9
    TRANSPARENT_WINDOW = 1 << 10
    FAST_PRESS = 1 << 11
    RAW_MOUSE_DISCREPANCY = 1 << 12
    RAW_KEYBOARD_DISCREPANCY = 1 << 13

IGNORE_HAX_FLAGS = Flags.CLEAN.value | Flags.INCORRECT_MOD.value