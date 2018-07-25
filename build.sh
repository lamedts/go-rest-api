#!/bin/bash

set -e

WORKING_DIR=$( dirname "${BASH_SOURCE[0]}" )
BUILD_DIR=$WORKING_DIR/build
DISTRO=linux_x64
COMPRESSION_METHOD="xz"
#DISTRO=darwin
VERSION=""
BUILD_DATE=$(date "+%Y-%m-%d@%H:%M:%S@%Z")
#BUILD_TAGS="''"
BUILD_TAGS="'EVDEV'"
NUM_RE='^[0-9]+$'

#CHECK_TEST=1

if [ $# == 0 ]; then
    echo -n 'VERSION_DERIVATIVE [optional]: ' 
    read VERSION_DERIVATIVE
	echo -n "VERSION_MAJOR: "
    read VERSION_MAJOR;
	echo -n "VERSION_MINOR: "
	read VERSION_MINOR;
    echo -n 'VERSION_RELEASE: ' 
    read VERSION_RELEASE
    echo -n 'VERSION_BUILD [optional, 4 chars]: ' 
    read VERSION_BUILD
    echo -n 'Compression: (xz) ' 
    read COMPRESSION_METHOD
elif [ $# == 4 ]; then
	VERSION_MAJOR=$1
	VERSION_MINOR=$2
	VERSION_RELEASE=$3
    COMPRESSION_METHOD="xz"
else
    echo "Error: Not enough parameters, enter either 3 or none parameters"
    exit 1;
fi

COMPRESSION_METHOD=${COMPRESSION_METHOD:-xz}

if [ ! -z $CHECK_TEST ]; then
    echo "Running tests before building..."
    set +e

    test_result=99
    if [[ "$DISTRO" == "linux_x64" ]]; then
        env GOOS=linux GOARCH=amd64 go test -cover --tags $BUILD_TAGS $WORKING_DIR/...
    elif [[ "$DISTRO" == "darwin" ]]; then
        env GOOS=darwin GOARCH=amd64 go test -cover --tags $BUILD_TAGS $WORKING_DIR/...
    fi
    test_result=$?

    if [ $test_result -ne 0 ]; then
        echo "Test Failed. Abort building. Please check before compiling."
        exit 2
    else
        echo "Test Passed"
    fi

    set -e
else
    echo "Skipping tests"
fi

if ! [[ $VERSION_MAJOR =~ $NUM_RE  ]] || ! [[ $VERSION_MINOR =~ $NUM_RE  ]]  || ! [[ $VERSION_RELEASE =~ $NUM_RE  ]] ; then
   echo "[error] some field is not a number" >&2; exit 1
elif [ "$COMPRESSION_METHOD" != "xz" ] && [ "$COMPRESSION_METHOD" != "gz" ]  ; then
    echo "[error] Wrong COMPRESSION_METHOD" >&2; exit 1
elif [ ${#VERSION_BUILD} -eq 4 ] ; then
    VERSION="${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_RELEASE}.${VERSION_BUILD}"
elif [ ${#VERSION_BUILD} -eq 0 ] ; then
    VERSION="${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_RELEASE}"
else
    echo "[error] Wrong VERSION_BUILD" >&2; exit 1
fi

VERSION_DERIVATIVE=`echo "$VERSION_DERIVATIVE" | tr '[:upper:]' '[:lower:]'`
VERSION=`echo "$VERSION" | tr '[:upper:]' '[:lower:]'`
VERSION_BUILD=`echo "$VERSION_BUILD" | tr '[:upper:]' '[:lower:]'`
if [ ${#VERSION_DERIVATIVE} -ne 0 ] ; then
    VERSION="${VERSION_DERIVATIVE}_${VERSION}"
fi
echo "Building Version: $VERSION for Distribution: $DISTRO"

DISTRO_BUILD_DIR=$BUILD_DIR/$DISTRO
if [ ! -d $BUILD_DIR ]; then
	mkdir -p $BUILD_DIR
fi
if [ -d $DISTRO_BUILD_DIR ]; then
	echo "Removing previous build..."
	rm -rf $DISTRO_BUILD_DIR
fi
mkdir -p $DISTRO_BUILD_DIR

set +e
build_result=99
if [[ "$DISTRO" == "linux_x64" ]]; then
	env GOOS=linux GOARCH=amd64 go build --tags $BUILD_TAGS -ldflags "-X main.BUILD=$BUILD_TAGS -X main.BUILD_DATE=$BUILD_DATE -X main.VERSION_MAJOR=$VERSION_MAJOR -X main.VERSION_MINOR=$VERSION_MINOR -X main.VERSION_RELEASE=$VERSION_RELEASE -X main.VERSION_DERIVATIVE=$VERSION_DERIVATIVE" -o $DISTRO_BUILD_DIR/go-rest-api -v go-rest-api/cmd/yay 2>&2
elif [[ "$DISTRO" == "darwin" ]]; then
    env GOOS=darwin GOARCH=amd64 go build --tags $BUILD_TAGS -ldflags "-X main.BUILD=$BUILD_TAGS -X main.BUILD_DATE=$BUILD_DATE -X main.VERSION_MAJOR=$VERSION_MAJOR -X main.VERSION_MINOR=$VERSION_MINOR -X main.VERSION_RELEASE=$VERSION_RELEASE -X main.VERSION_DERIVATIVE=$VERSION_DERIVATIVE" -o $DISTRO_BUILD_DIR/go-rest-api -v go-rest-api/cmd/yay 2>&2
fi
build_result=$?
if [ $build_result -ne 0 ]; then
        echo "Build Failed."
        exit 3
else
    echo "Finished building"
fi

set -e

cp $GOPATH/src/go-rest-api/config/config.yaml $DISTRO_BUILD_DIR/

echo "Build finished, Now packing into a zipped file"

DISTRO_ZIPPED_FILE=$BUILD_DIR/$DISTRO.tar
if [ -f $DISTRO_ZIPPED_FILE ]; then
    rm $DISTRO_ZIPPED_FILE
fi

if [ "$COMPRESSION_METHOD" = "gz" ] ; then
	cd $BUILD_DIR; tar -zcf yay_$VERSION.tar.gz $DISTRO/* ; cd $WORKING_DIR
    echo "Created zip file yay_$VERSION.tar.gz from $BUILD_DIR/$DISTRO"
else
	cd $BUILD_DIR; XZ_OPT=-e9 tar -cJf yay_$VERSION.tar.xz $DISTRO/* ; cd $WORKING_DIR
    echo "Created zip file yay_$VERSION.tar.xz from $BUILD_DIR/$DISTRO"
fi

# for docker
DOCKER_BUILD_DIR=../docker/build
if [ -d $DOCKER_BUILD_DIR ]; then
	echo "Removing previous build in docker folder"
	rm -rf $DOCKER_BUILD_DIR
fi
mkdir $DOCKER_BUILD_DIR
cp -r * $DOCKER_BUILD_DIR