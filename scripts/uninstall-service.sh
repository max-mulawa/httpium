#set -e
pgrep httpium
if [ $? -eq 0 ]; then
    echo "stopping service"
    systemctl stop httpium
    systemctl disable httpium     
fi

if [ -f "/etc/systemd/system/httpium.service" ]; then
    rm -f /etc/systemd/system/httpium.service
    rm -f /usr/sbin/httpium
    systemctl daemon-reload
    echo "httpium service removed"
fi 

