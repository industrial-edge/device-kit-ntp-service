#!/bin/bash
# Copyright © Siemens 2020 - 2025. ALL RIGHTS RESERVED.
# Licensed under the MIT license
# See LICENSE file in the top-level directory

chmod 711 /usr/bin/ntpservice 
chmod 640 /lib/systemd/system/dm-ntp.service
systemctl daemon-reload
systemctl enable dm-ntp.service
systemctl restart dm-ntp.service