/*
 * @Author              : Lihang
 * @Email               : lihang818@foxmail.com
 * @Date                : 2025-10-10 14:53:38
 * @LastEditTime        : 2025-10-16 09:51:40
 * @Description         :
 */
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func runConsumer() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("test-consumer-group"),
		consumer.WithNameServer([]string{
			"122.248.211.86:9876",
		}),
	)

	err := c.Subscribe("TestTopic", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			fmt.Printf("收到消息: %s\n", string(msg.Body))
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(err)
	}

	_ = c.Start()
	fmt.Println("消费者已启动，等待消息...")
	time.Sleep(time.Hour)
	_ = c.Shutdown()
}
