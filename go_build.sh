#!/bin/bash

cd ~/uam-power-backend

gos=(
    "id_route_main.go"
    "task_route_main.go"
    "data_route_main.go"
    "lane_route_main.go"
    "transfer_main.go"
)

for cmd in "${gos[@]}"; do
    echo "Building: $cmd"
    eval "go build $cmd"
done

echo "All go are built!"
