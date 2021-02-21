curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"device":"entrance","type":1}' \
  http://rpihome:3000/detect

rpihome -d entrance -t 1