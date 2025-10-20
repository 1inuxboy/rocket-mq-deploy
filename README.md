<!--
 * @Author              : Lihang
 * @Date                : 2025-10-20 15:52:35
 * @Description         : 
 * @Email               : lihang818@foxmail.com
 * @LastEditTime        : 2025-10-20 16:06:12
-->
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
./bin/mqbroker -c /opt/rocketmq-4.9.7/conf/acl.conf
```
