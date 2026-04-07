package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "fleet.events"
	QueueName    = "geofence_alerts"
	RoutingKey   = "geofence.entry"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(amqpURL string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := ch.QueueBind(
		q.Name,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	log.Println("RabbitMQ exchange and queue declared successfully")

	return &Publisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *Publisher) Publish(payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		ExchangeName,
		RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) StartWorker() error {
	msgs, err := p.channel.Consume(
		QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("RabbitMQ worker received geofence event: %s", string(msg.Body))
		}
	}()

	log.Println("RabbitMQ worker started")
	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
