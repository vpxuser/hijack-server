#!/bin/bash
# ============================================================
# 主控脚本，用于控制WIFI热点和劫持服务的启动和关闭
# ============================================================

# 引入同目录下的其他脚本
source "./hotspot/deps.sh"

# 引入子目录的脚本
# ============================================================
# 主程序入口
# ============================================================
case "$1" in
  up|start)
    cd "./hotspot" || exit 1
    bash "./main.sh" "$1" || exit 1
    cd "../" || exit 1
    sudo chmod +x hijack_amd64_linux
    msg "启动劫持服务..."
    nohup ./hijack_amd64_linux &
    ;;
  down|stop)
    cd "./hotspot" || exit 1
    bash "./main.sh" "$1" || exit 1
    cd "../" || exit 1
    sudo pkill -f hijack_amd64_linux || true
    msg "劫持服务已关闭"
    ;;
  status)
    cd "./hotspot" || exit 1
    bash "./main.sh" "$1" || exit 1
    ;;
  *)
    echo "用法: $0 {up|down|status}"
    ;;
esac
