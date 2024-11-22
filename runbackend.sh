#!/bin/bash

# 获取当前日期和时间（格式：YYYY-MM-DD_HH-MM-SS）
timestamp=$(date +"%Y-%m-%d_%H-%M-%S")

# 创建日志文件夹
log_dir=" ~/uam-power-backend/logs/$timestamp"
mkdir -p "$log_dir"

# 定义要运行的 nohup 命令数组
commands=(
    "nohup ~/uam-power-backend/id_route_main > $log_dir/id_route_main.log 2>&1 &"
    "nohup ~/uam-power-backend/task_route_main > $log_dir/task_route_main.log 2>&1 &"
    "nohup ~/uam-power-backend/data_route_main > $log_dir/data_route_main.log 2>&1 &"
    "nohup ~/uam-power-backend/transfer_main > $log_dir/transfer_main.log 2>&1 &"
)

# 循环执行每个 nohup 命令
for cmd in "${commands[@]}"; do
    echo "Executing: $cmd"
    eval $cmd
done

echo "All programs are running. Logs are saved in $log_dir"
