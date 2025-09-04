#!/bin/bash
# ============================================================
# WiFi 热点主脚本
# ============================================================

# 导入模块
source ./deps.sh
source ./config.sh
source ./start.sh
source ./stop.sh
source ./status.sh

# ============================================================
# 主程序入口
# ============================================================
case "$1" in
  up|start)
    start_hotspot
    ;;
  down|stop)
    stop_hotspot
    ;;
  status)
    show_status
    ;;
  *)
    echo "用法: $0 {up|down|status}"
    ;;
esac
