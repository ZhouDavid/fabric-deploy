#!/bin/bash

ip=$1
path=$2
# docker-compose check
docker-compose -f ${path}/docker-compose/orderers/${ip}.yaml up -d
