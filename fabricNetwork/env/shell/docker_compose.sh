#!/bin/sh
type=$1
binFile=$2
distPath="/usr/local/bin"

# 安装
function install() {
  cp $binFile $distPath
  which docker-compose
  if [ "$?" -ne 0 ]; then
    echo "docker-compose tool not found. exiting"
  fi
  echo 'docker-compose successed'
}
install
