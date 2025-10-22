## 部署方式

> 使用RocketMQ的4.9.7版本

**默认nameserver已经使用docker完成了部署**

1. 二进制
2. docker部署

## 二进制部署 broker

```bash
# java version
java -version 

java version "1.8.0_151"
Java(TM) SE Runtime Environment (build 1.8.0_151-b12)
Java HotSpot(TM) 64-Bit Server VM (build 25.151-b12, mixed mode)

wget https://github.com/apache/rocketmq/archive/refs/tags/rocketmq-all-4.9.7.tar.gz
tar xf rocketmq-all-4.9.7.tar.gz
cd rocketmq-rocketmq-all-4.9.7/
mvn -Prelease-all -DskipTests -Dspotbugs.skip=true clean install -U
cp -r distribution/target/rocketmq-4.9.7.tar.gz /opt/
tar xf rocketmq-4.9.7.tar.gz
ll
```

### 配置文件

详见 `deploy/rocketmq-conf`

### 启动和测试

```bash
ln -sf /opt/rocketmq-4.9.7/bin/* /usr/local/bin/
mqbroker -c /opt/rocketmq-4.9.7/conf/broker.conf 

# 设置环境变量（可选，不设置则使用默认地址）
export ROCKETMQ_NAMESERVER="<ip>:9876"

# 运行测试程序
go run main.go
```

#### 环境变量说明

- `ROCKETMQ_NAMESERVER`: RocketMQ NameServer 地址
  - 支持单个地址: `<ip>:9876`
  - 支持多个地址（逗号分隔）: `192.168.1.1:9876,192.168.1.2:9876`
  - 支持多个地址（分号分隔）: `192.168.1.1:9876;192.168.1.2:9876`
  - 如果不设置此环境变量，默认使用: `<ip>:9876`

> 使用systemctl管理，参考 `deploy/systemctl`

```bash
cat <<EOF > /etc/systemd/system/rocketmq-broker.service
[Unit]
Description=Apache RocketMQ Broker (Single Mode)
After=network.target docker.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/rocketmq-4.9.7
ExecStart=/opt/rocketmq-4.9.7/bin/mqbroker -c /opt/rocketmq-4.9.7/conf/broker.conf
Restart=on-failure
RestartSec=10
StandardOutput=append:/data/logs/broker.log
StandardError=append:/data/logs/broker.log

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable rocketmq-broker --now
```

## Docker 部署

待开发测试完成后，更新此文档