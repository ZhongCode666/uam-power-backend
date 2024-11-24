#!/bin/bash

echo "Build and run!"

echo "Creating kafka topic!"

/opt/kafka_2.12-3.8.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create  --if-not-exists --topic AircraftData10Area --partitions 25   --replication-factor 1
echo "Create data topic successfully!"

/opt/kafka_2.12-3.8.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create  --if-not-exists --topic AircraftEvent10Area --partitions 5   --replication-factor 1
echo "Create event topic successfully!"

cd ~/uam-power-backend

echo "CD to working dir!"

gos=(
    "id_route_main"
    "task_route_main"
    "data_route_main"
    "lane_route_main"
    "transfer_main"
)

for cmd in "${gos[@]}"; do
    echo "Building: $cmd.go"
    eval "go build $cmd.go"
done

echo "All go are built!"


# 获取当前日期和时间（格式：YYYY-MM-DD_HH-MM-SS）
timestamp=$(date +"%Y-%m-%d_%H-%M-%S")

# 创建日志文件夹
log_dir="$HOME/uam-power-backend/logs/$timestamp"
mkdir -p "$log_dir"


# 循环执行每个 nohup 命令
for cmd in "${gos[@]}"; do
    echo "Executing: $cmd"
    eval "nohup $HOME/uam-power-backend/$cmd > $log_dir/$cmd.log 2>&1 &"
done

echo "All programs are running. Logs are saved in $log_dir"
