#!/bin/bash

set -e

build() {
    if [[ ! -d deployments/bin ]];then
        mkdir deployments/bin
    fi

    dirname=./app/interface/$2/cmd
    if [ -d $dirname ];then
		for f in ${dirname}/main.go; do \
		    if [ -f ${f} ];then \
		        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o deployments/bin/$1/$2_$1/$2_$1 ${dirname}
                echo build over: $1_$2; \
            fi \
		done \
	fi
}

buildall() {
    #web
    build web comment
    build web post
    build web user

    #srv
    build srv comment
    build srv identify
    build srv post
    build srv user
}

case $1 in
    all) echo "全部build"
    buildall
    ;;
    build) echo "build:"$2,$3
    if [[ -z $2 || -z $3 ]];then
    echo "param error"
    exit 2
    fi
    build $2 $3
    ;;
    *)
    echo "build error"
    exit 2
    ;;
esac
