#!/bin/bash

WORKSPACE=$(
    cd $(dirname $0)/
    pwd
)
cd $WORKSPACE

function start() {
    # replace null to nil
    sed -i 's/null//g' ../docker/docker-compose-ehl.yaml
    docker-compose -f ../docker/docker-compose-ehl.yaml up -d
}

start
