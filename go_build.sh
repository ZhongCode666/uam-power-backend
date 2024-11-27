#!/bin/bash

cd ~/uam-power-backend

apis_gos=(
    "id_route_main"
    "task_route_main"
    "upload_route_main"
    "lane_route_main"
    "area_route_main"
    "receive_route_main"
)

transfer_gos=(
  "transfer_to_mysql"
  "transfer_to_redis"
)

echo "Starting to build apis..."

for cmd in "${apis_gos[@]}"; do
    echo "Building: main_services/apis/$cmd.go"
    eval "go build main_services/apis/$cmd.go"
done

echo "Starting to build transfer..."

for cmd in "${transfer_gos[@]}"; do
    echo "Building: main_services/transfer/$cmd.go"
    eval "go build main_services/transfer/$cmd.go"
done

echo "All go are built!"
