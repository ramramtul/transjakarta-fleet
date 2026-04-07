package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	mqttlib "github.com/eclipse/paho.mqtt.golang"

	"transjakarta-fleet/internal/config"
	"transjakarta-fleet/internal/database"
	"transjakarta-fleet/internal/handler"
	mqttsubscriber "transjakarta-fleet/internal/mqtt"
	"transjakarta-fleet/internal/rabbitmq"
	"transjakarta-fleet/internal/repository"
	"transjakarta-fleet/internal/router"
	"transjakarta-fleet/internal/service"
)

func main() {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		log.Fatal("POSTGRES_DSN is required")
	}

	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		log.Fatal("MQTT_BROKER is required")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL is required")
	}

	var db *sql.DB
	err := config.Retry(10, 3*time.Second, "PostgreSQL", func() error {
		var err error
		db, err = database.NewPostgres(postgresDSN)
		return err
	})
	if err != nil {
		log.Fatal("failed to connect PostgreSQL after retries: ", err)
	}

	var rabbitPub *rabbitmq.Publisher
	err = config.Retry(10, 3*time.Second, "RabbitMQ", func() error {
		var err error
		rabbitPub, err = rabbitmq.NewPublisher(rabbitURL)
		return err
	})
	if err != nil {
		log.Fatal("failed to connect RabbitMQ after retries: ", err)
	}
	defer rabbitPub.Close()

	if err := rabbitPub.StartWorker(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewLocationRepository(db)
	locationService := service.NewLocationService(repo, rabbitPub)

	var mqttClient mqttlib.Client
	err = config.Retry(10, 3*time.Second, "MQTT", func() error {
		opts := mqttlib.NewClientOptions()
		opts.AddBroker(mqttBroker)
		opts.SetClientID("fleet-backend-subscriber")

		mqttClient = mqttlib.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		return nil
	})
	if err != nil {
		log.Fatal("failed to connect MQTT after retries: ", err)
	}

	log.Printf("Connected to MQTT broker: %s", mqttBroker)

	subscriber := mqttsubscriber.NewSubscriber(mqttClient, locationService)
	if err := subscriber.Start(); err != nil {
		log.Fatal(err)
	}

	log.Println("MQTT subscriber started")

	vehicleHandler := handler.NewVehicleHandler(locationService)
	r := router.SetupRouter(vehicleHandler)

	log.Println("HTTP server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
