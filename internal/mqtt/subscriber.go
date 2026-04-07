package mqtt

import (
	"context"
	"encoding/json"
	"log"

	"transjakarta-fleet/internal/domain"
	"transjakarta-fleet/internal/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber struct {
	client  mqtt.Client
	service *service.LocationService
}

func NewSubscriber(client mqtt.Client, service *service.LocationService) *Subscriber {
	return &Subscriber{
		client:  client,
		service: service,
	}
}

func (s *Subscriber) Start() error {
	topic := "/fleet/vehicle/+/location"

	token := s.client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("MQTT received topic=%s payload=%s", msg.Topic(), string(msg.Payload()))

		var payload domain.LocationMessage
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Printf("invalid json payload: %v", err)
			return
		}

		if err := s.service.ProcessLocation(context.Background(), payload); err != nil {
			log.Printf("failed to process location: %v", err)
			return
		}

		log.Printf("location processed successfully for vehicle=%s", payload.VehicleID)
	})

	token.Wait()
	return token.Error()
}
