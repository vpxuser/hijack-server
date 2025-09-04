#!/bin/bash
# ============================================================
# WiFi 热点依赖安装和检查模块
# ============================================================

msg() {
  echo -e "[+] $1"
}

err() {
  echo -e "[!] $1" >&2
}

check_cmd() {
  command -v "$1" >/dev/null 2>&1
}

install_pkg() {
  local pkg=$1
  if check_cmd apt-get; then
    sudo apt-get update
    sudo apt-get install -y "$pkg"
  elif check_cmd dnf; then
    sudo dnf install -y "$pkg"
  elif check_cmd yum; then
    sudo yum install -y "$pkg"
  else
    err "未找到支持的包管理器 (apt/dnf/yum)"
    exit 1
  fi
}

install_deps() {
  msg "检查 hostapd..."
  if ! check_cmd hostapd; then
    msg "hostapd 未安装，正在安装..."
    install_pkg hostapd
  else
    msg "hostapd 已安装"
  fi

  msg "检查 dnsmasq..."
  if ! check_cmd dnsmasq; then
    msg "dnsmasq 未安装，正在安装..."
    install_pkg dnsmasq
  else
    msg "dnsmasq 已安装"
  fi

  msg "检查 iw..."
  if ! check_cmd iw; then
    msg "iw 未安装，正在安装..."
    install_pkg iw
  fi

  msg "检查 iproute2..."
  if ! check_cmd ip; then
    msg "iproute2 未安装，正在安装..."
    install_pkg iproute2
  fi

  msg "检查 iptables/nftables..."
  if ! check_cmd iptables && ! check_cmd nft; then
    msg "iptables/nftables 未安装，正在安装 iptables..."
    install_pkg iptables
  fi
}
