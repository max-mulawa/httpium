#!/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./httpium ./cmd/server

sudo systemctl stop httpium
sudo systemctl disable httpium

echo "Copy service definition"
sudo cp ./httpium.service /etc/systemd/system/

echo "Copy service binaries"
sudo mv ./httpium /usr/sbin/

echo "Reload services definitions"
sudo systemctl daemon-reload
echo "Activate httpium service"
sudo systemctl enable httpium
echo "Start httpium service"
sudo systemctl start httpium

systemctl status httpium
