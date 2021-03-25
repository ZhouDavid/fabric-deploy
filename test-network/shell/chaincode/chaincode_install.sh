#!/bin/bash
function join { local IFS="$1"; shift; echo "$*"; }
function getOrgDomain() {
    address=$1
    local orgDomain=(${address//:/ })
    orgDomain=${orgDomain[0]}
    orgDomain=(${orgDomain//./ })
    join "." ${orgDomain[@]:1}
}
function getDomain() {
    address=$1
    local items=(${address//:/ })
    echo ${items[0]}
}
function getPort() {
    address=$1
    local items=(${address//:/ })
    echo ${items[1]}
}

chainCodeType="go"
dockerName=$1 # doker name 默认 cli
chaincodePath=$2 # path in docker
chaincodeName=$3
chaincodeVersion=$4
peerAddress=$5
ordererAddress=$6
channelName=$7

ordererOrgDomain=$(getOrgDomain $ordererAddress)
ordererDomain=$(getDomain $ordererAddress)
ordererPort=$(getPort $ordererAddress)

echo "ordererOrgDomain=$ordererOrgDomain"
echo "ordererDomain=$ordererDomain"
echo "ordererPort=$ordererPort"

peerOrgDomain=$(getOrgDomain $peerAddress)
peerDomain=$(getDomain $peerAddress)
peerPort=$(getPort $peerAddress)

function goInstallPackage() {
    # go mod
    docker exec cli /bin/sh -c "cd chaincode/go/basic;go mod vendor"
    # package
    docker exec cli peer lifecycle chaincode package "/opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode/go/$chaincodeName.tar.gz" --path $chaincodePath --lang golang --label basic_1.0    
    # install
    res=`docker exec cli peer lifecycle chaincode install /opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode/go/$chaincodeName.tar.gz 2>&1 1>/dev/null`
    # package_id
    items=(${res// / })
    package_id=${items[25]}
    # approve definition
    docker exec cli peer lifecycle chaincode approveformyorg -o "$ordererDomain:$ordererPort" --ordererTLSHostnameOverride "$ordererDomain" --channelID $channelName --name $chaincodeName --version $chaincodeVersion --package-id $package_id --sequence 1 --tls --cafile "/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$ordererOrgDomain/orderers/$ordererDomain/msp/tlscacerts/tlsca.$ordererOrgDomain-cert.pem"
    # check commit readiness
    docker exec cli peer lifecycle chaincode checkcommitreadiness --channelID $channelName --name $chaincodeName --version $chaincodeVersion --sequence 1
    # commit
    docker exec cli peer lifecycle chaincode commit -o "$ordererDomain:$ordererPort" --ordererTLSHostnameOverride "$ordererDomain" --tls --cafile "/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$ordererOrgDomain/orderers/$ordererDomain/msp/tlscacerts/tlsca.$ordererOrgDomain-cert.pem" --channelID $channelName --name $chaincodeName --peerAddresses $peerAddress --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/$peerOrgDomain/peers/$peerDomain/tls/ca.crt --version $chaincodeVersion --sequence 1
    # initial invoke
    docker exec cli peer chaincode invoke -o "$ordererDomain:$ordererPort" --ordererTLSHostnameOverride "$ordererDomain" --tls --cafile "/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$ordererOrgDomain/orderers/$ordererDomain/msp/tlscacerts/tlsca.$ordererOrgDomain-cert.pem" -C $channelName -n $chaincodeName --peerAddresses $peerAddress --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/$peerOrgDomain/peers/$peerDomain/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'
}

if [[ ${chainCodeType} == "go" ]]; then
    goInstallPackage
else
    echo "Chaincode type is not supported yet."
fi
