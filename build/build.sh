#!/bin/bash
../build/build_gateway.sh
../build/build_news.sh

docker build -t 2charm/sql ../database