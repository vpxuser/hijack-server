#!/bin/bash
# ============================================================
# 停止 WiFi 热点模块
# ============================================================

stop_hotspot() {
  msg "停止 hostapd 和 dnsmasq..."
  sudo pkill hostapd || true
  sudo pkill dnsmasq || true

  msg "清理 NAT 规则..."
  if [ "$USE_NFT" = true ] && check_cmd nft; then
    sudo nft flush ruleset
  else
    sudo iptables -t nat -F
  fi

  msg "热点已停止"
}
