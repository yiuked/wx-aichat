#!/bin/bash

# shellcheck disable=SC1068
source=" ./ghapi ./docker-compose.yaml"
zipfile="./ghapi.zip"
remotedir="/data/wx-aichat"
host="$1"
make="install"

# build code
#git pull
set -x
make "$make"
set +x

# clear old zip file
if [ -f "$zipfile" ]; then
  rm -f $zipfile
fi

# zip deploy files
zip -q -r $zipfile $source
echo "info:zip finished"

scp -r $zipfile root@$host:$remotedir

ssh root@$host <<eeooff
cd $remotedir
rm -rf $source __MACOSX
unzip $zipfile
rm -f $zipfile

# backup docker logs
#docker-compose down
docker compose restart
eeooff

# clear file
rm -f ghapi
rm -f $zipfile
