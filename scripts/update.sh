#!/bin/bash

RED='\033[0;31m'         #  ${RED}
GREEN='\033[0;32m'      #  ${GREEN}
NORMAL='\033[0m'      #  ${NORMAL}

# Обновляем ОС.

echo -e "${GREEN}Updating system${NORMAL}"
apt-get -yqq update
apt-get -yqq dist-upgrade
echo -e "${GREEN}finish updating${NORMAL}"

# Обновляем приложение из репозитория.

echo -e "${GREEN}Update app${NORMAL}"
cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home
git push origin master
echo -e "${GREEN}Finish updating app${NORMAL}"

# Компилируем бинарник приложения.

echo -e "${GREEN}Compile app${NORMAL}"
/usr/local/go/bin/go build -o detector -v ./cmd/detector/main.go && /usr/local/go/bin/go build -o daemon -v ./cmd/daemon/main.go
echo -e "${GREEN}Finish compiling app{NORMAL}"
