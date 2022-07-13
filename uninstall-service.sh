sudo systemctl stop httpium
sudo systemctl disable httpium
sudo rm -f /etc/systemd/system/httpium.service
sudo rm -f /usr/sbin/httpium
sudo systemctl daemon-reload