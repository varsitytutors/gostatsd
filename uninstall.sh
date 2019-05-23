#!/bin/bash
systemctl stop gostatsd
systemctl disable gostatsd
rm -rf /opt/gostatsd
rm -rf /etc/systemd/system/gostatsd.service
