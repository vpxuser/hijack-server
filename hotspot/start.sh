#!/bin/bash
# ============================================================
# 启动 WiFi 热点模块（支持 iptables 和 nftables TCP 全端口转发）
# ============================================================

start_hotspot() {
  install_deps
  ask_params
  gen_config

  # 检查无线网卡是否支持 AP 模式
  if ! iw list 2>/dev/null | grep -q AP; then
    err "当前无线网卡不支持 AP 模式"
    exit 1
  fi

  msg "配置 $WLAN_IF 为 AP 网关..."
  sudo ip link set $WLAN_IF down || true
  sudo ip addr flush dev $WLAN_IF || true
  sudo ip addr add 10.42.0.1/24 dev $WLAN_IF
  sudo ip link set $WLAN_IF up

  msg "启动 dnsmasq..."
  sudo pkill dnsmasq || true
  sudo dnsmasq -C "$DNSMASQ_CONF"

  msg "启动 hostapd..."
  sudo pkill hostapd || true
  sudo hostapd -B "$HOSTAPD_CONF"

  msg "启用 NAT 转发..."
  sudo sysctl -w net.ipv4.ip_forward=1

  if [ "$USE_NFT" = true ] && check_cmd nft; then
      # 创建 NAT 表和链
      sudo nft add table ip nat 2>/dev/null || true
      sudo nft "add chain ip nat POSTROUTING { type nat hook postrouting priority 100 ; }" 2>/dev/null || true
      sudo nft "add chain ip nat PREROUTING { type nat hook prerouting priority 0 ; }" 2>/dev/null || true

      # MASQUERADE
      sudo nft add rule ip nat POSTROUTING oifname "$WAN_IF" masquerade

      # 创建 filter 表和 FORWARD 链
      sudo nft add table ip filter 2>/dev/null || true
      sudo nft "add chain ip filter FORWARD { type filter hook forward priority 0 ; }" 2>/dev/null || true

      # FORWARD 规则
      sudo nft add rule ip filter FORWARD iifname "$WAN_IF" oifname "$WLAN_IF" ct state related,established accept
      sudo nft add rule ip filter FORWARD iifname "$WLAN_IF" oifname "$WAN_IF" accept

      if [ -n "$TARGET_PORT" ]; then
          msg "设置 nftables 转发，将 $WLAN_IF 的 TCP 所有端口流量转发到 0.0.0.0:$TARGET_PORT"
          # 转发所有 TCP 端口到指定 IP:PORT
          # 转发 TCP 所有端口 1-65535 到本地 8084
          sudo nft add rule ip nat PREROUTING iif wlan0 tcp dport 1-65535 redirect to $TARGET_PORT
      else
          msg "未输入目标，热点将以普通模式启动，不进行流量转发"
      fi
  else
      # iptables 规则
      sudo iptables -t nat -A POSTROUTING -o "$WAN_IF" -j MASQUERADE
      sudo iptables -A FORWARD -i "$WAN_IF" -o "$WLAN_IF" -m state --state RELATED,ESTABLISHED -j ACCEPT
      sudo iptables -A FORWARD -i "$WLAN_IF" -o "$WAN_IF" -j ACCEPT

      if [ -n "$TARGET_PORT" ]; then
          msg "设置 iptables 转发，将 $WLAN_IF 的 TCP 所有端口流量转发到 0.0.0.0:$TARGET_PORT"
          # 转发所有 TCP 流量到本地指定端口
          sudo iptables -t nat -A PREROUTING -i $WLAN_IF -p tcp -j REDIRECT --to-port $TARGET_PORT
      else
          msg "未输入目标，热点将以普通模式启动，不进行流量转发"
      fi
  fi

  msg "热点已启动: SSID=$SSID 密码=$PASSPHRASE"
}
