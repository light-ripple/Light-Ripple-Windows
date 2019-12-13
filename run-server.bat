start cmd /C "C: & cd C:\Users\tanvi\Pictures\nginx-1.16.1 & title NGINX & nginx.exe"
start cmd /C "C: & cd "C:\Users\tanvi\Pictures\Redis" & title Redis & redis-server.exe"
timeout 3
start cmd /C "C: & cd C:\Users\tanvi\Pictures\pep.py & title Bancho & python pep.py"
start cmd /C "C: & cd C:\Users\tanvi\Pictures\lets & title LETS & python lets.py"
start cmd /C "C: & cd C:\Users\tanvi\Pictures\avatar-server & title Avatar-Server & python avatarserver.py"
start bash.exe -c "cd /mnt/d/osu\!\ server/Frontend && ./frontend"
start bash.exe -c "cd /mnt/d/osu\!\ server/api && ./API"