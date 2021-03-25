#!/bin/sh

WORKSPACE=$(
    cd $(dirname $0)/
    pwd
)

rm -rf ${WORKSPACE}/../organizations
rm -rf ${WORKSPACE}/../ipmap.txt
rm -rf ${WORKSPACE}/../configtx/*
rm -rf ${WORKSPACE}/../docker-compose/orderers/*
rm -rf ${WORKSPACE}/../docker-compose/peer/*
rm -rf ${WORKSPACE}/../system-genesis-block/*
rm -rf ${WORKSPACE}/../*.yaml
