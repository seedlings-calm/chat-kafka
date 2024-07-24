#!/usr/bin/env bash

docker-compose -f ./internal/interfaces/kafka/docker-compose.yml up -d

docker ps

# docker-compose -f ./internal/interfaces/kafka/docker-compose.yml logs kafka-0
# docker-compose -f ./internal/interfaces/kafka/docker-compose.yml logs kafka-1
# docker-compose -f ./internal/interfaces/kafka/docker-compose.yml logs kafka-2
# docker-compose -f ./internal/interfaces/kafka/docker-compose.yml logs kafdrop