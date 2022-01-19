#!/bin/bash

chmod 711 /usr/bin/ntpservice 
chmod 640 /lib/systemd/system/dm-ntp.service
systemctl daemon-reload
systemctl enable dm-ntp.service
systemctl restart dm-ntp.service