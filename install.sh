#!/bin/bash
mkdir -p /opt/gostatsd
systemctl stop gostatsd
cp gostatsd /opt/gostatsd
cp gostatsd.toml /opt/gostatsd
cp gostatsd.service /etc/systemd/system/gostatsd.service
systemctl enable gostatsd
systemctl start gostatsd
