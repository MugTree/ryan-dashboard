#!/bin/bash

DOMAIN_NAME="somesite.co.uk"
DB_NAME="database.sqlite"
USER="portable@portable.cpp"
REMOTE_DIRECTORY="/home/portable/somesite.co.uk/www/"
BACKUP_DIR="db_backups/"
LOCAL_DIRECTORY="/Users/me/Developer/go-projects/www-template/data/archive/"
SERVICE_NAME="somesite.service"
TIME_SLUG=$(date +%F_%T)

while true; do
    read -p "Stop www.${DOMAIN_NAME}?: " yn
    case $yn in
    [Yy]*)

        # stop the service and copy the existing db to a timestamped file
        ssh -n $USER "sudo systemctl stop $SERVICE_NAME && 
        mkdir ${REMOTE_DIRECTORY}${BACKUP_DIR}${TIME_SLUG} &&  
        cp ${REMOTE_DIRECTORY}${DB_NAME} ${REMOTE_DIRECTORY}${BACKUP_DIR}${TIME_SLUG}/${DB_NAME} &&
        cp ${REMOTE_DIRECTORY}${DB_NAME}-* ${REMOTE_DIRECTORY}${BACKUP_DIR}${TIME_SLUG} &&
        echo "\""db rescued ${TIME_SLUG}"\"" > ${REMOTE_DIRECTORY}rescued.txt && exit;"
        echo "stopped $SERVICE_NAME ... and copied to ${REMOTE_DIRECTORY}${BACKUP_DIR}${TIME_SLUG}/"

        # copy the copied files down from the remote
        mkdir ${LOCAL_DIRECTORY}${TIME_SLUG}
        scp ${USER}:${REMOTE_DIRECTORY}${BACKUP_DIR}${TIME_SLUG}/${DB_NAME} ${LOCAL_DIRECTORY}${TIME_SLUG}/${DB_NAME}
        scp ${USER}:${REMOTE_DIRECTORY}${BACKUP_DIR}${TIME_SLUG}/${DB_NAME}-* ${LOCAL_DIRECTORY}${TIME_SLUG}
        echo "securely copied files down to ${LOCAL_DIRECTORY}${TIME_SLUG}"

        # restart the server
        ssh -n $USER "sudo systemctl start $SERVICE_NAME;"
        echo "restarted $SERVICE_NAME ..."

        echo "opening rescued db  ${LOCAL_DIRECTORY}${TIME_SLUG}/${DB_NAME} ..."
        # open the db using datagrip
        open -na "DataGrip.app" ${LOCAL_DIRECTORY}${TIME_SLUG}/${DB_NAME}

        break
        ;;
    [Nn]*) exit ;;
    *) echo "Please answer y/n." ;;
    esac
done
