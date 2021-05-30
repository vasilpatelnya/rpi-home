#!/bin/bash

RED='\033[0;31m'         #  ${RED}
GREEN='\033[0;32m'      #  ${GREEN}
NORMAL='\033[0m'      #  ${NORMAL}

# Обновляем ОС.

echo -e "${GREEN}Updating system${NORMAL}"
apt-get -yqq update
apt-get -yqq dist-upgrade
echo -e "${GREEN}finish updating${NORMAL}"

# Устанавливаем docker и docker-compose

curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
sudo usermod -aG docker $USER
rm get-docker.sh

sudo apt install -y libffi-dev libssl-dev python3-dev
sudo apt install -y python3 python3-pip
sudo pip3 install docker-compose
