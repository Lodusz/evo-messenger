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
		broker = "localhost:9092"
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  "chat-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Service{writer: w, reader: r}
}

func (s *Service) ProduceMessage(msg repository.Message) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = s.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: bytes,
		},
	)
	if err != nil {
		log.Printf("Ошибка отправки в Kafka: %v", err)
	}
	return err
}

func (s *Service) ConsumeMessages(msgChan chan<- repository.Message) {
	go func() {
		log.Println("Воркер Kafka Consumer успешно запущен")
		for {
			m, err := s.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Ошибка чтения из Kafka: %v", err)
				continue
			}

			var msg repository.Message
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("Ошибка парсинга JSON из Kafka: %v", err)
				continue
			}

			msgChan <- msg
		}
	}()
}
