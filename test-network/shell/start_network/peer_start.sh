#!/bin/bash
ip=$1
path=$2
# docker-compose check
docker-compose -f ${path}/docker-compose/peers/${ip}.yaml up -d
