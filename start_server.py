from pynput.keyboard import Controller
import subprocess
from time import sleep
from keyboard import press

keyboard = Controller()

def peppy():
    subprocess.Popen([r"C:\Program Files\ConEmu\ConEmu64.exe"])

    sleep(0.8)

    for characters in "cd C:/Users/tanvi/Pictures/pep.py":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(0.4)

    for characters in "python pep.py":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(1.2)

def lets():
    subprocess.Popen([r"C:\Program Files\ConEmu\ConEmu64.exe"])

    sleep(0.8)

    for characters in "cd C:/Users/tanvi/Pictures/lets":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(0.4)

    for characters in "python lets.py":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(1)

def avatar_server():
    subprocess.Popen([r"C:\Program Files\ConEmu\ConEmu64.exe"])

    sleep(0.8)

    for characters in "cd C:/Users/tanvi/Pictures/avatar-server":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(0.4)

    for characters in "python avatarserver.py":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(1)

def redis():
    subprocess.Popen([r"C:/Users/tanvi/Pictures/Redis/redis-server.exe"])

    sleep(0.8)

def nginx():
    subprocess.Popen([r"C:\Windows\System32\cmd.exe"])

    sleep(0.8)

    for characters in "cd C:/Users/tanvi/Pictures/nginx-1.16.1":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(0.4)

    for characters in "nginx":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

def frontend():
    subprocess.Popen([r"C:\Windows\System32\bash.exe"])

    sleep(0.8)

    for characters in "cd /mnt/d/'osu! server'/Frontend":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(0.8)

    for characters in "./frontend":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

def api():
    subprocess.Popen([r"C:\Windows\System32\bash.exe"])

    sleep(0.8)

    for characters in "cd /mnt/d/'osu! server'/api":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

    sleep(0.8)

    for characters in "./API":
        keyboard.press(characters)
        keyboard.release(characters)
        sleep(0.1)
    press('enter')

def main():
    redis()
    peppy()
    lets()
    avatar_server()
    frontend()
    api()
    nginx()

if __name__ == '__main__':
    main()
