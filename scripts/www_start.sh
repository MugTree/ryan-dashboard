#!/bin/bash
DOMAIN_NAME="somesite.co.uk"
USER="portable@portable.cpp"
SERVICE_NAME="somesite.service"

while true; do
    read -p "Start www.${DOMAIN_NAME}?: " yn
    case $yn in
    [Yy]*)
        ssh -n $USER "sudo systemctl start $SERVICE_NAME && exit"
        echo "started $SERVICE_NAME ..."
        break
        ;;
    [Nn]*) exit ;;
    *) echo "Please answer y/n." ;;
    esac
done
