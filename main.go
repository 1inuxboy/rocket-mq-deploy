/*
 * @Author              : Lihang
 * @Email               : lihang818@foxmail.com
 * @Date                : 2025-10-10 15:47:43
 * @LastEditTime        : 2025-10-17 10:10:38
 * @Description         :
 */
package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	fmt.Println("=== RocketMQ ACL 认证测试 ===")

	// 测试场景1：使用正确的管理员账号
	test_acl_success()

	// 测试场景2：不提供任何账号信息
	test_acl_no_credentials()

	// 测试场景3：使用错误的账号信息
	test_acl_wrong_credentials()

	// 测试场景4：使用非管理员账号（根据plain-acl.conf，该账号对topicA和topic_plan_status_update只有DENY权限）
	test_acl_normal_user()

	select {}
}

// 测试场景1：使用正确的管理员账号
func test_acl_success() {
	fmt.Println("\n[测试场景1] 使用正确的管理员账号")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-admin-group"),
		producer.WithNameServer([]string{"122.248.211.86:9876"}),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: "rocketmq2",
			SecretKey: "12345678",
		}),
	)
	if err != nil {
		fmt.Printf("创建生产者失败: %s\n", err)
		return
	}

	// Start the producer before sending messages
	err = p.Start()
	if err != nil {
		fmt.Printf("启动生产者失败: %s\n", err)
		return
	}

	defer p.Shutdown()

	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello RocketMQ with ACL Admin!"),
	})
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", result.String())
	}
}

// 测试场景2：不提供任何账号信息
func test_acl_no_credentials() {
	fmt.Println("\n[测试场景2] 不提供任何账号信息")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-no-cred-group"),
		producer.WithNameServer([]string{"122.248.211.86:9876"}),
		// 故意不设置credentials
	)
	if err != nil {
		fmt.Printf("创建生产者失败: %s\n", err)
		return
	}

	err = p.Start()
	if err != nil {
		fmt.Printf("启动生产者失败: %s\n", err)
		return
	}

	defer p.Shutdown()

	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello RocketMQ without ACL!"),
	})
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", result.String())
	}
}

// 测试场景3：使用错误的账号信息
func test_acl_wrong_credentials() {
	fmt.Println("\n[测试场景3] 使用错误的账号信息")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-wrong-cred-group"),
		producer.WithNameServer([]string{"122.248.211.86:9876"}),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: "wrongUser",
			SecretKey: "wrongPassword",
		}),
	)
	if err != nil {
		fmt.Printf("创建生产者失败: %s\n", err)
		return
	}

	err = p.Start()
	if err != nil {
		fmt.Printf("启动生产者失败: %s\n", err)
		return
	}

	defer p.Shutdown()

	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello RocketMQ with Wrong ACL!"),
	})
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", result.String())
	}
}

// 测试场景4：使用非管理员账号访问受限资源
func test_acl_normal_user() {
	fmt.Println("\n[测试场景4] 使用非管理员账号访问受限资源")
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("test-normal-group"),
		producer.WithNameServer([]string{"122.248.211.86:9876"}),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: "RocketMQ",
			SecretKey: "12345678",
		}),
	)
	if err != nil {
		fmt.Printf("创建生产者失败: %s\n", err)
		return
	}

	err = p.Start()
	if err != nil {
		fmt.Printf("启动生产者失败: %s\n", err)
		return
	}

	defer p.Shutdown()

	// 根据plain-acl.conf，RocketMQ用户对topicA和topic_plan_status_update只有DENY权限
	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "topicA",
		Body:  []byte("Hello from normal user to denied topic!"),
	})
	if err != nil {
		fmt.Printf("发送到topicA失败(预期行为): %s\n", err)
	} else {
		fmt.Printf("发送到topicA成功(非预期行为): %s\n", result.String())
	}

	// 测试发送到topic_plan_status_update（也应该被拒绝）
	result, err = p.SendSync(context.Background(), &primitive.Message{
		Topic: "topic_plan_status_update",
		Body:  []byte("Hello from normal user to denied topic_plan_status_update!"),
	})
	if err != nil {
		fmt.Printf("发送到topic_plan_status_update失败(预期行为): %s\n", err)
	} else {
		fmt.Printf("发送到topic_plan_status_update成功(非预期行为): %s\n", result.String())
	}

	// 尝试发送到有权限的topicB
	result, err = p.SendSync(context.Background(), &primitive.Message{
		Topic: "topicB",
		Body:  []byte("Hello from normal user to allowed topic!"),
	})
	if err != nil {
		fmt.Printf("发送到topicB失败(非预期行为): %s\n", err)
	} else {
		fmt.Printf("发送到topicB成功(预期行为): %s\n", result.String())
	}
}
