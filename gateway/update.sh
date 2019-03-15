#!/bin/bash
docker volume prune -f
docker rm -f api_gateway
docker network create api_gateway
#Set env variables
export ADDR=:443
export KEY="`cat ../news_api.key`"

docker pull 2charm/spectrum_gateway
docker run -d --network service_network --name api_gateway \
-p 443:443 \
-e ADDR=$ADDR \
-e KEY=$KEY \
2charm/spectrum_gateway