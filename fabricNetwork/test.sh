#!/bin/bash
# this script is only used for testing fabric network locally
root=$1

cd $root
echo $PWD
sudo ./cleanup.sh
echo "loading config..."
./fabricNetwork loadConfig --config $PWD/test-network/networkconfig.json
echo "start network locally..."
sudo ./fabricNetwork startNetwork --dPath $PWD/test-network --sPath $PWD/test-network
echo "creating channel..."
sleep 5
sudo ./fabricNetwork createChannel --dPath $PWD/test-network --channel-name mychannel

echo "installing basic chaincode..."
./fabricNetwork installChaincode --ccPath $PWD/test-network/chaincode/go/basic --ccName basic --ccVersion 1.0 --channelName mychannel --hosts peer0.org1.example.com