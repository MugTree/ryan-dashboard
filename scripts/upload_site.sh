#!/bin/bash
DOMAIN_NAME="www.somesite.co.uk"
PROJECT_NAME="somesite web application"
CURRENT_DATE=$(date "+%F_%T")
LOCAL_DIRECTORY="/Users/me/Developer/go-projects/ryan-dashboard"
REMOTE_DIRECTORY="/home/portable/somesite.co.uk/www"
USER="portable@portable.cpp"
HAS_UPLOADED=0
SERVICE_NAME="somesite.service"
MAKE_CMD=production-build-www
BINARY=somesite.www.amd64

while true; do
    read -p "Upload latest build of ${PROJECT_NAME} to ${DOMAIN_NAME}?: " yn
    case $yn in
    [Yy]*)
        cd $LOCAL_DIRECTORY && make $MAKE_CMD
        ssh -n $USER "sudo systemctl stop $SERVICE_NAME && exit"
        echo "stopped $SERVICE_NAME ..."
        echo "copying latest files ..."
        rsync -rv $LOCAL_DIRECTORY/bin/$BINARY $USER:$REMOTE_DIRECTORY/bin
        ssh -n $USER " echo uploaded! && sudo systemctl start $SERVICE_NAME &&  exit"
        echo "starting $SERVICE_NAME, $POSTGRES_SERVICE_NAME ..."
        HAS_UPLOADED=1
        break
        ;;
    [Nn]*) exit ;;
    *) echo "Please answer y/n." ;;
    esac
done
