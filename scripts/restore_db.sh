#!/bin/bash

DOMAIN_NAME="somesite.co.uk"
DB_NAME="database.sqlite"
USER="portable@portable.cpp"
REMOTE_DIRECTORY="/home/portable/www-template/www/"
BACKUP_DIR="db_backups/"
SERVICE_NAME="somesite.service"
TIME_SLUG=$(date +%F_%T)

# script takes a directory path as an argument...
# checks that the path exists and checks that the file we want to upload also exists
# eg...

# ./scripts/restore_db.sh /Users/me/Developer/go-projects/notzero/data

# strip an ending slash as could cause issues
# eg. $1="/data/" to "/data"
dir=${1%/}

if [ -d $dir ]; then
    echo "$dir exists"

    if find $dir -maxdepth 1 -type f -name "${DB_NAME}" | grep -q "."; then
        echo "$DB_NAME exists"

        while true; do
            read -p "Stop www.${DOMAIN_NAME} and restore the db from ${dir}?: " yn
            case $yn in
            [Yy]*)

                # stop the service and run some copy commands on the box
                ssh -n $USER "sudo systemctl stop $SERVICE_NAME && exit;"

                # copy the copied files up to the remote - these will overwrite existing files
                scp $dir/$DB_NAME ${USER}:${REMOTE_DIRECTORY}

                # copy up any wal files etc
                scp $dir/${DB_NAME}-* ${USER}:${REMOTE_DIRECTORY}

                # restart the server
                ssh -n $USER "cd ${REMOTE_DIRECTORY} && echo "\""db restored at ${TIME_SLUG}"\"" &&  sudo systemctl start $SERVICE_NAME && exit;"
                echo "restarted $SERVICE_NAME ..."

                break
                ;;
            [Nn]*) exit ;;
            *) echo "Please answer y/n." ;;
            esac
        done

    else
        echo "$DB_NAME doesn't exist"
    fi
else
    echo "$1 doesn't exist"
fi
