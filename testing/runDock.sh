#!/bin/bash

docker run \
  -it \
  --rm \
  --net=host \
  -v `pwd`/haproxy/:/var/run/haproxy/ \
  --name haproxy-syntax-check \
  haproxy haproxy -f /var/run/haproxy/haproxy.cfg
