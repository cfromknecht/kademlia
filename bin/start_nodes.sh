#!/bin/bash

MAIN_PATH=run/main.go

PORT=6001

MAX_NODES=50
COUNTER=0
while [ $COUNTER -lt $MAX_NODES ]; do
  go run $MAIN_PATH --port $PORT --first-ip 127.0.0.1:6000 --first-id $1 &
  let COUNTER=COUNTER+1
  let PORT=PORT+1
  sleep 1
done
