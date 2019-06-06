#!/bin/bash
../build/build_news.sh
docker push 2charm/news_service
scp -i "~/downloads/privatekey.pem" ../assets/news_api.key ec2-user@ec2-3-18-100-210.us-east-2.compute.amazonaws.com:~/
ssh -i "~/downloads/privatekey.pem" ec2-user@ec2-3-18-100-210.us-east-2.compute.amazonaws.com 'bash -s' < update_news.sh