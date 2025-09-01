# 劫持工具V2使用说明

## 服务介绍

**本劫持工具有以下三个服务：**

- **劫持服务**：端口可自定义，在主配置文件config/config.yml上修改即可，默认端口为8084
- **抓包服务**：端口写死为8080，不可修改，使用抓包服务时请先安装config/cacert.crt和config/cacert.key
- **报告查看服务**：端口写死为8081，访问8081端口即可在线查看报告

## 配置文件说明

### 主配置文件config/config.yml

```yaml
#控制台打印的日志等级
#当日志等级设置为 4 时
#控制台只会打印 ≤4 的日志内容
#1：fatal，2：error，3：warn，4：info，5：debug
#当为debug等级时，控制台会打印完整的请求和响应快照
#默认日志等级为info
logLevel: 4
#跳过ssl劫持失败的地址
#在第一次劫持失败后，后续再遇到这个失败的地址会不进行劫持
#如果这个地址存在被多个http client调用的情况，可能会存在漏测，若漏测，可设置为false
#默认开启跳过
skip: true
#报告生成开关，默认开启
report: true
#劫持证书路径
cert: config/hijack.crt
#劫持证书的私钥路径
key: config/hijack.key
#劫持服务监听的地址
host: 0.0.0.0
#劫持服务监听的端口
port: 8084
#需要删除的响应头清单，如某些安全头和缓存头会影响劫持效果，需要删除
headers: config/headers.yml
#自定义规则，可自定义替换响应体中链接的正则
qrcode: config/qrcode.yml
#黄色图片链接
imageURL: https://img95.699pic.com/xsj/0x/3k/lt.jpg!/fh/300
#自定义规则替换的黄色网站链接
qrcodeURL: https://www.baidu.com
#替换响应体图片的文件
image: config/hijack.png
#替换响应体html的文件
html: config/hijack.html
#需要进行劫持的域名清单
targets: config/target.list
#报告生成的模板文件，不要修改
template: config/template.html
```

### 响应体配置文件config/headers.yml

**可向其中添加你想删除的响应头字段，添加格式如下**：

```yaml
#按照yaml格式添加配置
#响应头分组名称:
# - header1
# - header2
#假设你要删除名为X-Request-ID、X-Forward-For的响应头，可以这样写
deleteHeaderGroup:
  - X-Request-ID
  - X-Forward-For
```

### 自定义正则配置文件config/qrcode.yml

**可向其中添加你想替换的响应体内的链接，添加格式如下**:

```yaml
#按照yaml格式添加配置
# 正则规则名称: >-
#  具体规则
#假设你要添加劫持微信二维码规则，可以这样写
微信二维码: >-
  (?i)https?://work\.weixin\.qq\.com/ct/wcde[a-zA-Z0-9]+(?:[/?][^\s]*)?
```

### 劫持目标配置文件config/target.list

**可向其中添加你想劫持的地址，添加格式如下**：

```tex
#按照文本文件格式添加，每行一个地址，可以这样写
www.baidu.com
www.google.com
192.168.0.1
172.28.1.1
```

