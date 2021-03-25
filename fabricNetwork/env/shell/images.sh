#!/bin/bash

WORKSPACE=$(
    cd $(dirname $0)/
    pwd
)
cd $WORKSPACE

FABRIC_IMAGES=(baseos.tar ca.tar ccenv.tar orderer.tar peer.tar tools.tar)

for image in ${FABRIC_IMAGES[@]}; do
    echo "load"
    # docker pull ${DOCKER_NS}/$image:${ARCH}-${VERSION}
    docker load -i "../images/"${image}
done
