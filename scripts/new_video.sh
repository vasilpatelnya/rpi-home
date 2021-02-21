#!/bin/bash

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"device":"entrance","type":2}' \
  http://rpihome:3000/api/v1/motioneye