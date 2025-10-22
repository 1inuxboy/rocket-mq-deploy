/*
 * @Author              : Lihang
 * @Email               : lihang818@foxmail.com
 * @Date                : 2025-10-10 15:47:43
 * @LastEditTime        : 2025-10-22 14:25:20
 * @Description         :
 */
package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/sirupsen/logrus"
)

var (
	testResults       = make(map[string]bool) // 存储测试结果
	rocketmqNameServ  []string                // RocketMQ NameServer 地址列表
	adminCredentials  primitive.Credentials   // 管理员凭证
	readerCredentials primitive.Credentials   // 普通用户凭证
)

func main() {
	// 完全禁用RocketMQ客户端内部日志
	os.Setenv("rocketmq.client.logLevel", "OFF")
	os.Setenv("rocketmq.client.logRoot", "/dev/null")
	os.Setenv("rocketmq.client.logFileMaxSize", "0")
	os.Setenv("rocketmq.client.logFileMaxIndex", "0")
	os.Setenv("rocketmq.client.logUseSlf4j", "false")
	rlog.SetLogLevel("ERROR")

	// 初始化 RocketMQ NameServer 配置
	initRocketMQConfig()

	// 设置简洁的日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableColors:          false,
		DisableLevelTruncation: true,
	})

	// 设置日志级别为Info，但只显示我们自己的日志
	logrus.SetLevel(logrus.InfoLevel)

	// 禁用标准库日志
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	logrus.Info("=== RocketMQ ACL 认证测试 ===")

	// 测试场景1：使用正确的管理员账号
	test_acl_success()

	// 测试场景2：不提供任何账号信息
	test_acl_no_credentials()

	// 测试场景3：使用错误的账号信息
	test_acl_wrong_credentials()

	// 测试场景4：使用非管理员账号（根据plain-acl.conf，该账号对topicA和topic_plan_status_update只有DENY权限）
	test_acl_normal_user()

	logrus.Info("=== 所有测试完成 ===")
	printTestSummary()
	select {}
}

// 初始化 RocketMQ 配置
func initRocketMQConfig() {
	readerAccessKey := os.Getenv("ROCKETMQ_READER_ACCESS_KEY")
	readerSecretKey := os.Getenv("ROCKETMQ_READER_SECRET_KEY")
	adminAccessKey := os.Getenv("ROCKETMQ_ADMIN_ACCESS_KEY")
	adminSecretKey := os.Getenv("ROCKETMQ_ADMIN_SECRET_KEY")
	nameServEnv := os.Getenv("ROCKETMQ_NAMESERVER")
	if readerAccessKey == "" || readerSecretKey == "" || adminAccessKey == "" || adminSecretKey == "" || nameServEnv == "" {
		logrus.Errorf("环境变量 ROCKETMQ_READER_ACCESS_KEY 或 ROCKETMQ_READER_SECRET_KEY 或 ROCKETMQ_ADMIN_ACCESS_KEY 或 ROCKETMQ_ADMIN_SECRET_KEY 未设置")
		os.Exit(1)
	}

	// 支持多个地址，使用逗号分隔
	addresses := strings.Split(nameServEnv, ",")
	for _, addr := range addresses {
		addr = strings.TrimSpace(addr)
		if addr != "" {
			rocketmqNameServ = append(rocketmqNameServ, addr)
		}
	}

	readerCredentials = primitive.Credentials{
		AccessKey: readerAccessKey,
		SecretKey: readerSecretKey,
	}
	adminCredentials = primitive.Credentials{
		AccessKey: adminAccessKey,
		SecretKey: adminSecretKey,
	}
}

// 打印测试结果汇总
func printTestSummary() {
	logrus.Info("=== 测试结果汇总 ===")
	passedCount := 0
	totalCount := len(testResults)

	for testName, passed := range testResults {
		if passed {
			logrus.Infof("✅ %s: 通过", testName)
			passedCount++
		} else {
			logrus.Errorf("❌ %s: 失败", testName)
		}
	}

	logrus.Infof("=== 总计: %d/%d 个测试通过 ===", passedCount, totalCount)
}

// 测试场景1：使用正确的管理员账号
func test_acl_success() {
	logrus.Info("🔍 [测试场景1] 使用正确的管理员账号")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-admin-group"),
		producer.WithNameServer(rocketmqNameServ),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: adminCredentials.AccessKey,
			SecretKey: adminCredentials.SecretKey,
		}),
	)
	if err != nil {
		logrus.Errorf("❌ 创建生产者失败: %s", err)
		testResults["场景1-管理员账号"] = false
		return
	}

	err = p.Start()
	if err != nil {
		logrus.Errorf("❌ 启动生产者失败: %s", err)
		testResults["场景1-管理员账号"] = false
		return
	}

	defer p.Shutdown()

	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello RocketMQ with ACL Admin!"),
	})
	if err != nil {
		logrus.Errorf("❌ 发送失败: %s (不符合预期)", err)
		testResults["场景1-管理员账号"] = false
	} else {
		logrus.Infof("✅ 发送成功: %s (符合预期)", result.String())
		testResults["场景1-管理员账号"] = true
	}
}

// 测试场景2：不提供任何账号信息
func test_acl_no_credentials() {
	logrus.Info("🔍 [测试场景2] 不提供任何账号信息")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-no-cred-group"),
		producer.WithNameServer(rocketmqNameServ),
		// 故意不设置credentials
	)
	if err != nil {
		logrus.Errorf("❌ 创建生产者失败: %s", err)
		testResults["场景2-无账号信息"] = false
		return
	}

	err = p.Start()
	if err != nil {
		logrus.Errorf("❌ 启动生产者失败: %s", err)
		testResults["场景2-无账号信息"] = false
		return
	}

	defer p.Shutdown()

	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello RocketMQ without ACL!"),
	})
	if err != nil {
		logrus.Infof("✅ 发送失败: %s (符合预期)", err)
		testResults["场景2-无账号信息"] = true
	} else {
		logrus.Errorf("❌ 发送成功: %s (不符合预期)", result.String())
		testResults["场景2-无账号信息"] = false
	}
}

// 测试场景3：使用错误的账号信息
func test_acl_wrong_credentials() {
	logrus.Info("🔍 [测试场景3] 使用错误的账号信息")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-wrong-cred-group"),
		producer.WithNameServer(rocketmqNameServ),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: "wrongUser",
			SecretKey: "wrongPassword",
		}),
	)
	if err != nil {
		logrus.Errorf("❌ 创建生产者失败: %s", err)
		testResults["场景3-错误账号"] = false
		return
	}

	err = p.Start()
	if err != nil {
		logrus.Errorf("❌ 启动生产者失败: %s", err)
		testResults["场景3-错误账号"] = false
		return
	}

	defer p.Shutdown()

	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello RocketMQ with Wrong ACL!"),
	})
	if err != nil {
		logrus.Infof("✅ 发送失败: %s (符合预期)", err)
		testResults["场景3-错误账号"] = true
	} else {
		logrus.Errorf("❌ 发送成功: %s (不符合预期)", result.String())
		testResults["场景3-错误账号"] = false
	}
}

// 测试场景4：使用非管理员账号访问受限资源
func test_acl_normal_user() {
	logrus.Info("🔍 [测试场景4] 使用非管理员账号访问受限资源")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-normal-group"),
		producer.WithNameServer(rocketmqNameServ),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: readerCredentials.AccessKey,
			SecretKey: readerCredentials.SecretKey,
		}),
	)
	if err != nil {
		logrus.Errorf("❌ 创建生产者失败: %s", err)
		testResults["场景4-普通用户"] = false
		return
	}

	err = p.Start()
	if err != nil {
		logrus.Errorf("❌ 启动生产者失败: %s", err)
		testResults["场景4-普通用户"] = false
		return
	}

	defer p.Shutdown()

	// 根据plain-acl.conf，RocketMQ用户对topicA和topic_plan_status_update只有DENY权限
	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topicA",
		Body:  []byte("Hello from normal user to denied topic!"),
	})
	if err != nil {
		logrus.Infof("✅ 发送到topicA失败: %s (符合预期)", err)
		testResults["场景4-topicA权限"] = true
	} else {
		logrus.Errorf("❌ 发送到topicA成功: %s (不符合预期)", result.String())
		testResults["场景4-topicA权限"] = false
	}

	// 测试发送到topic_plan_status_update（也应该被拒绝）
	result, err = p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello from normal user to denied topic_plan_status_update!"),
	})
	if err != nil {
		logrus.Infof("✅ 发送到topic_plan_status_update失败: %s (符合预期)", err)
		testResults["场景4-topic_plan_status_update权限"] = true
	} else {
		logrus.Errorf("❌ 发送到topic_plan_status_update成功: %s (不符合预期)", result.String())
		testResults["场景4-topic_plan_status_update权限"] = false
	}

	// 尝试发送到有权限的topicB
	result, err = p.SendSync(context.Background(), &primitive.Message{
		Topic: "topicB",
		Body:  []byte("Hello from normal user to allowed topic!"),
	})
	if err != nil {
		logrus.Errorf("❌ 发送到topicB失败: %s (不符合预期)", err)
		testResults["场景4-topicB权限"] = false
	} else {
		logrus.Infof("✅ 发送到topicB成功: %s (符合预期)", result.String())
		testResults["场景4-topicB权限"] = true
	}
}
