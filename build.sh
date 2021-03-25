#!/bin/bash
cd $(dirname $0)
DIR=$(pwd)
tools="fabricNetwork"
cd $tools
mkdir -p target

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o target/fabricNetwork main.go

cp -r env target/
cp -r ${DIR}/test-network target/
# cp -r sampleconfig/*.json target/
cp -r ${DIR}/test-network/*.json target/
cd target
tar -zcvf fabricNetwork.tar.gz *
