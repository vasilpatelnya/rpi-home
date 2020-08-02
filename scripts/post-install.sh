#!/bin/bash

cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/detect.sh /var/lib/motioneye/detect.sh
cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/new_video.sh /var/lib/motioneye/new_video.sh

cd ..
/usr/local/go/bin/go build -o detector -v /home/pi/go/src/github.com/vasilpatelnya/rpi-home/cmd/detector/main.go && /usr/local/go/bin/go build -o daemon -v /home/pi/go/src/github.com/vasilpatelnya/rpi-home/cmd/daemon/main.go
/home/pi/go/src/github.com/vasilpatelnya/rpi-home/daemon &