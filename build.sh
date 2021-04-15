#!/bin/bash
cd $(dirname $0)
DIR=$(pwd)
#find . -type f -name "*.sh" -exec chmod 0755 {}\;
#find . -type f -name "*.sh" -exec dos2unix {} \+

mkdir -p target/test-network
pushd fabricNetwork
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../target/fabricNetwork main.go
popd
cp -r fabricNetwork/env/ target/
cp cleanup.sh target/
cp test.sh target/
cp -r ${DIR}/test-network/chaincode target/test-network/
cp -r ${DIR}/test-network/explorer target/test-network/
cp -r ${DIR}/test-network/shell target/test-network/
cp -r ${DIR}/test-network/networkconfig.json target/test-network/

pushd target
tar -zcvf fabricNetwork.tar.gz *
mv fabricNetwork.tar.gz ..
popd
rm -rf target
