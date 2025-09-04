#!/bin/bash
# ============================================================
# WiFi 热点配置交互模块（检查无线网卡，若无则退出）
# ============================================================

CONF_DIR=/etc/wifi-hotspot
HOSTAPD_CONF="$CONF_DIR/hostapd.conf"
DNSMASQ_CONF="$CONF_DIR/dnsmasq.conf"

ask_params() {
  # 检查无线驱动是否存在
  if ! iw dev >/dev/null 2>&1; then
    echo "检测不到无线驱动或 nl80211 接口，尝试加载常见内核模块..."
    sudo modprobe cfg80211 || true
    sudo modprobe mac80211 || true
    sudo modprobe ath9k || true  # 示例驱动
    sudo modprobe iwlwifi || true  # 示例驱动
  fi

  # 检测无线接口
  WLAN_LIST=$(iw dev 2>/dev/null | grep Interface | awk '{print $2}')
  if [ -z "$WLAN_LIST" ]; then
    echo "未检测到可用无线网卡，请确认无线网卡已插入且支持 AP 模式。"
    exit 1
  fi

  echo "可用无线网卡接口："
  echo "$WLAN_LIST"

  read -rp "请输入无线网卡接口 (默认 wlan0): " WLAN_IF
  WLAN_IF=${WLAN_IF:-wlan0}

  echo "可用外网接口："
  ip link show | awk -F: '/^[0-9]+: [^lo]/ {print $2}' | sed 's/ //g'

  read -rp "请输入外网接口 (默认 eth0): " WAN_IF
  WAN_IF=${WAN_IF:-eth0}

  read -rp "请输入无线 SSID (默认 HijackWIFI): " SSID
  SSID=${SSID:-HijackWIFI}

  read -rp "请输入热点密码 (默认 88888888): " PASSPHRASE
  PASSPHRASE=${PASSPHRASE:-88888888}

  read -rp "请选择无线模式 g=2.4GHz, a=5GHz (默认 g): " HW_MODE
  HW_MODE=${HW_MODE:-g}

  read -rp "请输入信道 (默认 6): " CHANNEL
  CHANNEL=${CHANNEL:-6}

  read -rp "是否使用 nftables (true/false, 默认 true): " USE_NFT
  USE_NFT=${USE_NFT:-true}

  # 可选流量转发 host:port
  read -rp "请输入劫持服务端口: " TARGET_PORT
}

gen_config() {
  sudo mkdir -p "$CONF_DIR"

  cat <<EOF | sudo tee "$HOSTAPD_CONF" >/dev/null
interface=$WLAN_IF
driver=nl80211
ssid=$SSID
hw_mode=$HW_MODE
channel=$CHANNEL
wmm_enabled=1
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase=$PASSPHRASE
wpa_key_mgmt=WPA-PSK
rsn_pairwise=CCMP
EOF

  cat <<EOF | sudo tee "$DNSMASQ_CONF" >/dev/null
interface=$WLAN_IF
dhcp-range=10.42.0.10,10.42.0.100,12h
dhcp-option=3,10.42.0.1
dhcp-option=6,8.8.8.8,8.8.4.4
EOF

  msg "配置文件已生成: $CONF_DIR"
}
