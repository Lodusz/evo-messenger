package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"chat-service/internal/repository"

	"github.com/segmentio/kafka-go"
)

const topic = "chat-messages"

type Service struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaService() *Service {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		log.Fatal("КРИТИЧЕСКАЯ ОШИБКА: KAFKA_BROKER не задан в окружении")
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "chat-group",
	})

	return &Service{writer: w, reader: r}
}

func (s *Service) ProduceMessage(msg repository.Message) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.writer.WriteMessages(context.Background(), kafka.Message{Value: bytes})
}

func (s *Service) ConsumeMessages(msgChan chan<- repository.Message) {
	go func() {
		for {
			m, err := s.reader.ReadMessage(context.Background())
			if err != nil {
				continue
			}
			var msg repository.Message
			if err := json.Unmarshal(m.Value, &msg); err == nil {
				msgChan <- msg
			}
		}
	}()
}
