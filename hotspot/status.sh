#!/bin/bash
# ============================================================
# WiFi 热点状态模块
# ============================================================

show_status() {
  ip addr show $WLAN_IF
  ip route show
  sudo pgrep -a hostapd || true
  sudo pgrep -a dnsmasq || true
}
