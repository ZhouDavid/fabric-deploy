#!/bin/sh

## go path ;/usr/local
installType=$1
tarFile=$2
distPath="/usr/local"
version="1.16"
url="https://studygolang.com/dl/golang/go${version}.linux-amd64.tar.gz"
# url="https://golang.org/dl/go${version}.linux-amd64.tar.gz"

if [ "${installType}" == "" ]; then
    echo 'default online'
    installType="online"
fi

function installOffline() {
    cp $tarFile $distPath
    tar -C $distPath -zxvf $tarFile
    echo 'export GOROOT=/usr/local/go' >>/etc/profile
    echo 'export GOPATH=/root/go' >>/etc/profile
    echo 'export PATH=$PATH:$GOROOT/bin' >>/etc/profile
    source /etc/profile
    echo 'go installed'
    go version
}

function installOnLine() {
    wget ${url}
    tar -C /usr/local -xzf go${version}.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >>/etc/profile
    source /etc/profile
    echo 'go installed'
    go version
}

if [ "${installType}" == "online" ]; then
    echo 'exec installOnLine '
    installOnLine
elif [ "${installType}" == "offline" ]; then
    echo 'exec installOffline'

    if [ "$tarFile" == "" ]; then
        exit 1
    fi

    installOffline
fi
