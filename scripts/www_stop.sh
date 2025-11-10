#!/bin/bash
DOMAIN_NAME="somesite.co.uk"
USER="portable@portable.cpp"
SERVICE_NAME="somesite.service"

while true; do
    read -p "Stop www.${DOMAIN_NAME}?: " yn
    case $yn in
    [Yy]*)
        ssh -n $USER "sudo systemctl stop $SERVICE_NAME && exit"
        echo "stopped $SERVICE_NAME ..."
        break
        ;;
    [Nn]*) exit ;;
    *) echo "Please answer y/n." ;;
    esac
done
