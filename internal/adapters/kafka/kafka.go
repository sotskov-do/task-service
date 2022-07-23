package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.com/g6834/team26/task/pkg/uuid"
)

type KafkaClient struct {
	KafkaConn *kafka.Conn
}

func New(url, topic string) (*KafkaClient, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", url, topic, 0)
	if err != nil {
		return nil, err
	}
	return &KafkaClient{KafkaConn: conn}, nil
}

func (KafkaClient *KafkaClient) ActionTask(ctx context.Context, u, t, v string) error {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	actionTaskReq := &KafkaAnalyticMessage{
		UUID:      u,
		Timestamp: time.Now().Unix(),
		Type:      t,
		Value:     v,
	}

	messageData, err := json.Marshal(actionTaskReq)
	if err != nil {
		return err
	}

	_, err = KafkaClient.KafkaConn.WriteMessages(
		kafka.Message{
			Key:   []byte(uuid.GenUUID()),
			Value: messageData,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (KafkaClient *KafkaClient) Stop(ctx context.Context) error {
	err := KafkaClient.KafkaConn.Close()
	if err != nil {
		return err
	}
	return nil
}

type KafkaAnalyticMessage struct {
	UUID      string `json:"uuid"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}
