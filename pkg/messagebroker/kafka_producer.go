package messagebroker

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type EventProducer struct {
	producer sarama.SyncProducer
}

func NewEventProducer(brokers []string) (*EventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll // Ensure log durability

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &EventProducer{producer: p}, nil
}

type TicketSoldEvent struct {
	OrderID  string `json:"order_id"`
	UserID   string `json:"user_id"`
	TicketID string `json:"ticket_id"`
	Email    string `json:"email"`
}

// PublishTicketSold publishes a message to the "ticket_sold_events" topic
func (p *EventProducer) PublishTicketSold(topic string, event TicketSoldEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}