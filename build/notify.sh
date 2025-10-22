#!/bin/bash

path=$1
filename=`/usr/bin/basename $path`
if [ ${filename:0:2} != ".#" ]
then
    buildtool-router ./ > ./router/router.go
fi
