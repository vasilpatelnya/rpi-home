#!/bin/bash

cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/detect.sh /var/lib/motioneye/detect.sh
cp /home/pi/go/src/github.com/vasilpatelnya/rpi-home/scripts/new_video.sh /var/lib/motioneye/new_video.sh

cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home
/usr/local/go/bin/go build -o detector -v ./cmd/detector/main.go && /usr/local/go/bin/go build -o daemon -v ./cmd/daemon/main.go

sed -e 's/exit 0/\n/' /etc/rc.local >> file
echo "cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home" >> file
echo "sudo ./daemon" >> file
echo "exit 0" >> file
cat file > /etc/rc.local
rm file