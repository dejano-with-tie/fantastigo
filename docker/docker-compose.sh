#!/bin/bash

readonly vehicle_id="$1"

docker_compose_down() {
   printf "\rSIGINT caught..."
   docker compose -f docker-compose.vehicle.yaml down
   sleep 5 # let the previous one shutdown
   docker compose -f docker-compose.yaml down
}
trap 'docker_compose_down' SIGINT

docker compose -f docker-compose.yaml up &
sleep 10 # let the previous one bootstrap

export THING_NAME="$vehicle_id"
THING_NAME="$vehicle_id" docker compose -p "$vehicle_id" -f docker-compose.vehicle.yaml config
THING_NAME="$vehicle_id" docker compose -p "$vehicle_id" -f docker-compose.vehicle.yaml up\
