[![Build Status](https://travis-ci.com/vasilpatelnya/rpi-home.svg?branch=master)](https://travis-ci.com/vasilpatelnya/rpi-home)
[![codecov](https://codecov.io/gh/vasilpatelnya/rpi-home/branch/master/graph/badge.svg)](https://codecov.io/gh/vasilpatelnya/rpi-home)
# rpi-home

Все приложение представляет собой:
 * демон, который отслеживает и обрабатывает записи о событиях в БД
 * cli-приложение, которое пишет события переданные в аргументах в БД

## Подготовка RPi.

1. Установить на SD-карту ОС Raspbian.
2. В корне раздела boot создать пустой файл с именем "ssh"
3. Подключить к модему RPi через LAN-кабель.
4. `ping raspberrypi.local` покажет ip (лучше чтобы в сети была одна RPi)
5. Подключаемся по ssh к RPi и меняем пароль пользователя pi.
6. [Устанавливаем motioneye:](https://groups.google.com/forum/#!topic/motioneye/wxdFOn2a28M)

   `sudo apt-get update`
   
   `sudo apt-get dist-upgrade`
   
   `sudo apt-get install libssl-dev libcurl4-openssl-dev libmariadbclient18 libpq5 mysql-common ffmpeg`
   
   `wget https://github.com/Motion-Project/motion/releases/download/release-4.0.1/pi_stretch_motion_4.0.1-1_armhf.deb`
   
   `sudo dpkg -i pi_stretch_motion_4.0.1-1_armhf.deb`
   
   `sudo pip install motioneye`
   
   `sudo mkdir -p /etc/motioneye`
   
   `sudo cp /usr/local/share/motioneye/extra/motioneye.conf.sample /etc/motioneye/motioneye.conf`
   
   `sudo mkdir -p /var/lib/motioneye`
   
   `sudo cp /usr/local/share/motioneye/extra/motioneye.systemd-unit-local /etc/systemd/system/motioneye.service`
   
   `sudo systemctl daemon-reload`
   
   `sudo systemctl enable motioneye`
   
   `sudo systemctl start motioneye`
   
   `sudo reboot`
   
9. [Устанавливаем MongoDB на RPi.](https://pimylifeup.com/mongodb-raspberry-pi/)
    
    `sudo apt update`
    
    `sudo apt upgrade`
    
    `sudo apt install mongodb`
    
    `sudo systemctl enable mongodb`
    
    `sudo systemctl start mongodb`
10. Устанавливаем Go на RPi.
    
    Скачиваем версию для RPi из списка на странице [https://golang.org/dl/](https://golang.org/dl/) `wget https://golang.org/dl/go1.14.4.linux-armv6l.tar.gz`

    Распаковываем: `sudo tar -C /usr/local -xzf go1.14.4.linux-armv6l.tar.gz`
    
    `export PATH=$PATH:/usr/local/go/bin`
    
    В конец файла .profile добавляем также "export PATH=$PATH:/usr/local/go/bin"
    
    Проверяем: `go version`
## Отправка сообщений на телеграм.

1. [Установить telegram-send.](https://github.com/rahiel/telegram-send) `sudo pip3 install telegram-send`
2. Обычно устанавливается в `/usr/local/bin`
3. Запустить `telegram-send --configure`, где ввести api telegram-бота, а после отослать ему код, сгенерированый приложением.

## Запись срабатывания в БД.

1. Скомпилировать приложение в корень (можно в любую директорию).
2. Прописать путь к приложению в настройках камеры в motioneye:
    
    `/home/pi/go/src/gitlab.com/vasilpatelnya/rpi-home/detector -device room -type 1`
    
    Где device это название камеры в motioneye, а тип со значением 1 это константа для детектирования камерой движения.

## Поместить скрипт detect.sh /var/lib/motioneye

    sudo chmod +x /var/lib/motioneye/detect.sh

    sudo chmod 777 /var/lib/motioneye/detect.sh
    
1. Прописать путь до скрипта в поле `Run A Command` в motioneye.

## Устанавливаем демон в автозагрузку например в /etc/rc.local

Перед `exit 0` добавляем:

    cd /home/pi/go/src/gitlab.com/vasilpatelnya/rpi-home
    sudo ./daemon
