# Устанавливаем motioneye.

apt-get update
apt-get dist-upgrade
apt-get install libssl-dev libcurl4-openssl-dev libmariadbclient18 libpq5 mysql-common ffmpeg
wget https://github.com/Motion-Project/motion/releases/download/release-4.0.1/pi_stretch_motion_4.0.1-1_armhf.deb
dpkg -i pi_stretch_motion_4.0.1-1_armhf.deb
pip install motioneye
mkdir -p /etc/motioneye
cp /usr/local/share/motioneye/extra/motioneye.conf.sample /etc/motioneye/motioneye.conf
mkdir -p /var/lib/motioneye
cp /usr/local/share/motioneye/extra/motioneye.systemd-unit-local /etc/systemd/system/motioneye.service
systemctl daemon-reload
systemctl enable motioneye
systemctl start motioneye

# Устанавливаем MongoDB.

apt update
apt upgrade
apt install mongodb
systemctl enable mongodb
systemctl start mongodb

# Устанавливаем Go.

wget https://golang.org/dl/go1.14.4.linux-armv6l.tar.gz
tar -C /usr/local -xzf go1.14.4.linux-armv6l.tar.gz
export PATH=$PATH:/usr/local/go/bin

# В конец файла .profile добавляем также "export PATH=$PATH:/usr/local/go/bin"

pip3 install telegram-send
telegram-send --configure

#Запись срабатывания в БД.
#Скомпилировать приложение в корень (можно в любую директорию).
#
#Прописать путь к приложению в настройках камеры в motioneye:
#
#/home/pi/go/src/github.com/vasilpatelnya/rpi-home/detector -device room -type 1
#
#Где device это название камеры в motioneye, а тип со значением 1 это константа для детектирования камерой движения.
#
#Поместить скрипт detect.sh /var/lib/motioneye
#sudo chmod +x /var/lib/motioneye/detect.sh
#
#sudo chmod 777 /var/lib/motioneye/detect.sh
#Прописать путь до скрипта в поле Run A Command в motioneye.
#Устанавливаем демон в автозагрузку например в /etc/rc.local
#Перед exit 0 добавляем:
#
#cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home
#sudo ./daemon