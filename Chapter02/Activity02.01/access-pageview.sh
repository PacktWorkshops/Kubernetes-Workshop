#!/usr/bin/env bash
# Sleep 1 second in each cycle.
SLEEP_INTERVAL=1

while true; do
    curl http://192.168.99.100:31953
    sleep $SLEEP_INTERVAL
done