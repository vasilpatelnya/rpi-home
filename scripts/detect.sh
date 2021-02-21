#cd /home/pi/go/src/github.com/vasilpatelnya/rpi-home && ./detector -device deviceName -type 1
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"device":"entrance","type":1}' \
  http://127.0.0.1:3000/detect