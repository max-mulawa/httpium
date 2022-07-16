#!/bin/bash
pgrep httpium
if [ $? -eq 0 ]; then
    echo "stopping service"
    systemctl stop httpium
    systemctl disable httpium     
fi

set -e

echo "Copy service definition"
cp ./httpium.service /etc/systemd/system/

echo "Copy service binaries"
mv ./bin/httpium /usr/sbin/

echo "Reload services definitions"
systemctl daemon-reload
echo "Activate httpium service"
systemctl enable httpium
echo "Start httpium service"
systemctl start httpium

systemctl status httpium
