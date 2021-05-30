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
cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home || exit
git pull origin master
echo -e "${GREEN}Finish updating app${NORMAL}"

# Компилируем бинарник приложения.

echo -e "${GREEN}Restart app${NORMAL}"
docker-compose stop
docker-compose build rpihome
docker-compose up -d
echo -e "${GREEN}App started${NORMAL}"
