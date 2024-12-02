#!/bin/bash

echo "Build and run!"

echo "Run python"
~/uam-gis-intelligence-service/runservice.sh
echo "run python successfully ~/uam-gis-intelligence-service/runservice.sh"

redis-cli FLUSHALL
echo "delete redis data successfully"
echo "Creating kafka topic!"

/opt/kafka_2.12-3.8.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create  --if-not-exists --topic AircraftData10Area --partitions 25   --replication-factor 1
echo "Create data topic successfully!"

/opt/kafka_2.12-3.8.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create  --if-not-exists --topic AircraftEvent10Area --partitions 5   --replication-factor 1
echo "Create event topic successfully!"

cd ~/uam-power-backend

echo "CD to working dir!"

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


# 获取当前日期和时间（格式：YYYY-MM-DD_HH-MM-SS）
timestamp=$(date +"%Y-%m-%d_%H-%M-%S")

# 创建日志文件夹
log_dir="$HOME/uam-power-backend/logs/$timestamp"
mkdir -p "$log_dir"


# 循环执行每个 nohup 命令
echo "Starting to execute apis..."
for cmd in "${apis_gos[@]}"; do
    echo "Executing: $cmd"
    eval "nohup $HOME/uam-power-backend/$cmd > $log_dir/$cmd.log 2>&1 &"
done

echo "Starting to execute transfer..."
for cmd in "${transfer_gos[@]}"; do
    echo "Executing: $cmd"
    eval "nohup $HOME/uam-power-backend/$cmd > $log_dir/$cmd.log 2>&1 &"
done


# 创建 latest_log 的快捷方式
ln -sfn "$log_dir" "$HOME/uam-power-backend/logs/latest_log"

echo "All programs are running. Logs are saved in $log_dir, quick view with 'cd $HOME/uam-power-backend/logs/latest_log'"
