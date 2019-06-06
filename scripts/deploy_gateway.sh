#!/bin/bash
../build/build_gateway.sh
docker push 2charm/gateway
ssh -i "~/downloads/privatekey.pem" ec2-user@ec2-3-18-100-210.us-east-2.compute.amazonaws.com 'bash -s' < update_gateway.sh