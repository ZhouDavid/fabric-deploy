channelName=$1
blockFile=/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$channelName.block

docker exec cli peer channel join -b $blockFile
