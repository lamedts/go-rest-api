#!/bin/bash

# sudo docker build -t go-rest-api docker/
# sudo docker run -p 8080:8080 go-rest-api 
# sudo docker run -d -it -p 27017:27017 --name=go-rest-api-1 go-rest-api /bin/sh
# sudo docker exec -i -t go-rest-api-1 /bin/bash
# docker-compose up -d

##
#  To build 
#
WORKING_DIR=$( dirname "${BASH_SOURCE[0]}" )
BUILD_DIR=$WORKING_DIR/build
DISTRO=linux_x64
DISTRO_BUILD_DIR=$BUILD_DIR/$DISTRO
if [ -d $DISTRO_BUILD_DIR ]; then
	echo -n "Previous build found, to rebuild(Y\n):"
    read REBUILD
else
    REBUILD=Y
fi

if [[ "$REBUILD" == "Y" ]]; then
    bash ./build.sh 1 0 0 gz
fi

##
#  start docker 
#
cd docker
OS=$(uname -s)
if [[ "$OS" == "Darwin" ]]; then
    docker-compose up --build
else
    sudo docker-compose up --build
fi