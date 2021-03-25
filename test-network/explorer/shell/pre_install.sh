#!/bin/sh
explorerUrl="hyperledger/explorer:1.1.4"
explorerDbUrl="hyperledger/explorer-db:1.1.4"
WORKSPACE=$(
    cd $(dirname $0)/
    pwd
)

createDbPath=${WORKSPACE}"/../db"
createDb="createdb.sh"
dockerImageOnline() {
    docker pull $explorerUrl
    docker pull $explorerDbUrl
}

postgressInstall() {
    yum -y install gcc gcc-c++ kernel-devel
    yum -y install https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm
    yum -y install postgresql96
    yum -y install postgresql96-server
    /usr/pgsql-9.6/bin/postgresql96-setup initdb
    systemctl enable postgresql-9.6
    systemctl start postgresql-9.6
}

jqInstall() {
    yum install wget -y
    wget http://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
    rpm -ivh epel-release-latest-7.noarch.rpm
    yum repolist
    yum install jq -y
}

nodeJsInstall() {
    curl -sL https://rpm.nodesource.com/setup_14.x | sudo bash -
    sudo yum install -y nodejs
}

start() {
    docker-compose -f ../docker/docker-compose-ehl.yaml up -d
}

postgresqlDBInit() {
    cd ${createDbPath}
    /usr/bin/bash $createDb
}

echo "docker pull"
dockerImageOnline

echo 'node js install'
nodeJsInstall

echo 'jq install'
jqInstall

echo 'postgress install'
postgressInstall

echo 'postgress db init'
postgresqlDBInit
