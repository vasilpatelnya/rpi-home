[![Build Status](https://travis-ci.com/vasilpatelnya/rpi-home.svg?branch=master)](https://travis-ci.com/vasilpatelnya/rpi-home)
[![codecov](https://codecov.io/gh/vasilpatelnya/rpi-home/branch/master/graph/badge.svg)](https://codecov.io/gh/vasilpatelnya/rpi-home)
# rpi-home

Все приложение представляет собой:
 * демон, который отслеживает и обрабатывает записи о событиях в БД
 * cli-приложение, которое пишет события переданные в аргументах в БД

![Схема приложения]("https://github.com/vasilpatelnya/rpi-home/blob/master/structure.svg")

## Подготовка RPi.

1. Установить на SD-карту ОС Raspbian.
2. В корне раздела boot создать пустой файл с именем "ssh"
3. Подключить к модему RPi через LAN-кабель.
4. `ping raspberrypi.local` покажет ip (лучше чтобы в сети была одна RPi)
5. Подключаемся по ssh к RPi и меняем пароль пользователя pi.
6. Клонируем репозиторий приложения `git clone https://github.com/vasilpatelnya/rpi-home`

## Автоматическая установка приложения (пока в тестовом режиме). Если проблемы с данным типом установке попробуйте пошаговую установку.

`cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home && make install`

## Запись срабатывания в БД.

1. Скомпилировать приложение в корень (можно в любую директорию).
2. Прописать путь к скриптам в настройках камеры в motioneye:

## Поместить скрипты detect.sh и new_video.sh в /var/lib/motioneye

    sudo cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/detect.sh /var/lib/motioneye/detect.sh

    sudo cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/new_video.sh /var/lib/motioneye/new_video.sh
    
    sudo chmod +x /var/lib/motioneye/detect.sh
    
    sudo chmod 777 /var/lib/motioneye/detect.sh
    
    sudo chmod +x /var/lib/motioneye/new_video.sh
    
    sudo chmod 777 /var/lib/motioneye/new_video.sh
    
1. Прописать путь до скриптов в поле `Run A Command` до detect.sh в поле `Run An End Command` до new_video.sh.

## Добавить в крон поднятие контейнеров при старте
