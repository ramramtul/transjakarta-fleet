package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		log.Fatal("MQTT_BROKER is required")
	}

	vehicleID := "B1234XYZ"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("mock-publisher")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Printf("Connected to MQTT broker: %s", broker)

	baseLat := -6.2088
	baseLng := 106.8456

	topic := "/fleet/vehicle/" + vehicleID + "/location"

	for {
		payload := LocationMessage{
			VehicleID: vehicleID,
			Latitude:  baseLat + (rand.Float64()-0.5)*0.0005,
			Longitude: baseLng + (rand.Float64()-0.5)*0.0005,
			Timestamp: time.Now().Unix(),
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("failed to marshal payload: %v", err)
			continue
		}

		token := client.Publish(topic, 1, false, body)
		token.Wait()

		if token.Error() != nil {
			log.Printf("failed to publish: %v", token.Error())
		} else {
			log.Printf("Published to %s: %s", topic, string(body))
		}

		time.Sleep(2 * time.Second)
	}
}
