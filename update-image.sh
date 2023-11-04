#!/bin/bash

echo "Tag: "
read tag

docker build -t qbqsystemsbot/socks5-relay:$tag .
docker login
docker push qbqsystemsbot/socks5-relay:$tag
docker logout
