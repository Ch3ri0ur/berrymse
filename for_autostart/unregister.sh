#!/bin/bash

# cp berrymse.service /etc/systemd/system/berrymse.service
# echo 'KERNEL=="video0", SYMLINK="video0", TAG+="systemd"' > /etc/udev/rules.d/webcam.rules

# systemctl daemon-reload

systemctl disable berrymse.service

rm /etc/udev/rules.d/webcam.rules

rm /etc/systemd/system/berrymse.service

systemctl daemon-reload