#!/usr/bin/env bash

docker-compose  up -d

docker ps

# docker-compose  logs kafka-0
# docker-compose  logs kafka-1
# docker-compose  logs kafka-2
# docker-compose  logs kafdrop