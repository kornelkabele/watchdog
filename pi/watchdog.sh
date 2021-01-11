#!/bin/bash
# Add to "crontab -e" following command as a user:
# @reboot /home/pi/watchdog/watchdog.sh &
# Make sure that both watchdog.sh and watchdog-pi have chmod u+x
cd /home/pi/watchdog
export $(cat .secrets | tr -d '\r' | xargs) && ./watchdog-pi
