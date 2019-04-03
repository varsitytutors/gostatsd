#!/bin/bash

mkdir -p /opt/gostatsd
systemctl stop gostatsd
cp gostatsd /opt/gostatsd
cp gostatsd.service /etc/systemd/system/gostatsd.service

dd_key=$(credstash -r $EC2_REGION get $ENVIRONMENT::datadog::api_key)
if [[ ! -z $dd_key ]]; then
  cat <<-EOF > /opt/gostatsd/gostatsd.toml
backends = 'newrelic datadog'

[datadog]
    api_key = "$dd_key"
EOF
else
  echo "Unable to find datadog api key - not adding datadog backend"
  echo "backends = 'newrelic'" > /opt/gostatsd/gostatsd.toml
fi

cat <<-EOF >> /opt/gostatsd/gostatsd.toml

[newrelic]
    address = "http://localhost:8001/v1/data"
    event-type = "StatsD"
EOF

systemctl enable gostatsd
systemctl start gostatsd
