/*
 * @Author              : Lihang
 * @Email               : lihang818@foxmail.com
 * @Date                : 2025-10-10 14:50:36
 * @LastEditTime        : 2025-10-16 09:52:14
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

func runProducer() {
	p, _ := rocketmq.NewProducer(
		producer.WithGroupName("test-producer-group"),
		producer.WithNameServer([]string{
			"122.248.211.86:9876",
		}),
	)
	err := p.Start()
	if err != nil {
		panic(err)
	}

	result, err := p.SendSync(context.Background(),
		&primitive.Message{
			Topic: "TestTopic",
			Body:  []byte("Hello RocketMQ from Go!"),
		},
	)
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", result.String())
	}

	p.Shutdown()
}
