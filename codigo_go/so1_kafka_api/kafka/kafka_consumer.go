package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)
func StartKafka()  {
	conf := kafka.ReaderConfig{
		Brokers: []string {"localhost:9092"},
		Topic: "topic_test2",
		MaxBytes: 10,
		SessionTimeout: 10 * time.Second,
	}

	reader := kafka.NewReader(conf)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Some error occured", err)
			continue
		}
		fmt.Println("Message is : ", string(m.Value))
	}

}