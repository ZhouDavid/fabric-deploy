#!/bin/bash
ordererAddress=$1 # for example, orderer0.orderer.example.com:7050
ordererDomain=$2
ordererOrgDomain=$3
channelName=$4
inputTxPath=/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$channelName.tx
outputBlockPath=/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$channelName.block
ordererCAPath=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$ordererOrgDomain/orderers/$ordererDomain/msp/tlscacerts/tlsca.$ordererOrgDomain-cert.pem

docker exec cli peer channel create -o $ordererAddress -c $channelName --ordererTLSHostnameOverride $ordererDomain -f $inputTxPath --outputBlock $outputBlockPath --tls --cafile $ordererCAPath
