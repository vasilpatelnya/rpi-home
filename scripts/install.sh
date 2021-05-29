#!/bin/bash

RED='\033[0;31m'         #  ${RED}
GREEN='\033[0;32m'      #  ${GREEN}
NORMAL='\033[0m'      #  ${NORMAL}

# Обновляем ОС.

echo -e "${GREEN}Updating system${NORMAL}"
apt-get -yqq update
apt-get -yqq dist-upgrade
echo -e "${GREEN}finish updating${NORMAL}"

# Устанавливаем motioneye.

echo -e "${GREEN}install motion eye${NORMAL}"
apt-get -yqq install libssl-dev libcurl4-openssl-dev libmariadbclient-dev libpq5 mysql-common ffmpeg libmicrohttpd12
wget https://github.com/Motion-Project/motion/releases/download/release-4.2.2/pi_buster_motion_4.2.2-1_armhf.deb
dpkg -i pi_buster_motion_4.2.2-1_armhf.deb
rm pi_buster_motion_4.2.2-1_armhf.deb
pip install motioneye
mkdir -p /etc/motioneye
cp /usr/local/share/motioneye/extra/motioneye.conf.sample /etc/motioneye/motioneye.conf
mkdir -p /var/lib/motioneye
cp /usr/local/share/motioneye/extra/motioneye.systemd-unit-local /etc/systemd/system/motioneye.service
systemctl daemon-reload
systemctl enable motioneye
systemctl start motioneye
echo -e "${GREEN}motioneye has been installed${NORMAL}"

# Устанавливаем MongoDB.

echo -e "${GREEN}install mongodb${NORMAL}"
apt -yqq update
apt -yqq upgrade
apt -yqq install mongodb
systemctl enable mongodb
systemctl start mongodb
echo -e ${GREEN}"mongodb has been installed${NORMAL}"

# Устанавливаем Go.

echo -e "${GREEN}golang installation${NORMAL}"
wget https://golang.org/dl/go1.14.4.linux-armv6l.tar.gz
tar -C /usr/local -xzf go1.14.4.linux-armv6l.tar.gz
echo "\nPATH=$PATH:/usr/local/go/bin\n" >> /home/pi/.profile
rm go1.14.4.linux-armv6l.tar.gz
echo -e "${GREEN}finish golang installation${NORMAL}"

# В конец файла .profile добавляем также "export PATH=$PATH:/usr/local/go/bin"

echo -e "${GREEN}telegram-send installation${NORMAL}"
pip3 install telegram-send
telegram-send --configure
echo -e "${GREEN}finish telegram-send installation${NORMAL}"

cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/detect.sh /var/lib/motioneye/detect.sh
cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/new_video.sh /var/lib/motioneye/new_video.sh

cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home
/usr/local/go/bin/go build -o rpihome -v ./cmd/rpihome/main.go

sed -e 's/exit 0/\n/' /etc/rc.local > file
echo "cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home" >> file
echo "sudo ./daemon -c configs/prod.env" >> file
echo "exit 0" >> file
cat file > /etc/rc.local
rm file