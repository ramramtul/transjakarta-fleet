# Transjakarta Fleet Management Backend

Backend service for vehicle fleet management system built with **Golang**, **MQTT**, **PostgreSQL**, **RabbitMQ**, and **Docker Compose**.

## Features

- Receive vehicle location data via **MQTT**
- Store vehicle location history in **PostgreSQL**
- Provide REST API for:
  - Latest vehicle location
  - Vehicle location history
- Publish **geofence events** to **RabbitMQ**
- Consume geofence events using background worker
- Run the full system with **Docker Compose**

---

## Tech Stack

- **Golang**
- **Gin**
- **PostgreSQL**
- **RabbitMQ**
- **Eclipse Mosquitto (MQTT Broker)**
- **Docker / Docker Compose**

---

## System Architecture

### Flow

1. Mock publisher sends vehicle GPS location to MQTT topic:

   `/fleet/vehicle/{vehicle_id}/location`

2. Backend subscribes to MQTT topic and validates incoming payload.

3. Valid location data is stored in PostgreSQL.

4. Backend checks whether vehicle enters a geofence area.

5. If vehicle enters geofence:
   - Publish event to RabbitMQ exchange `fleet.events`
   - Route event to queue `geofence_alerts`

6. Worker consumes messages from `geofence_alerts`.

7. REST API provides latest location and location history.

---

## Project Structure

```bash
cmd/
  app/            # Main backend app
  publisher/      # Mock MQTT publisher

internal/
  config/         # Retry helper / config utilities
  database/       # PostgreSQL connection
  domain/         # Data models
  handler/        # HTTP handlers
  mqtt/           # MQTT subscriber
  rabbitmq/       # RabbitMQ publisher + worker
  repository/     # Database access
  router/         # Gin router
  service/        # Business logic

migrations/       # SQL migration files
```


## Project Structure

Make sure you have installed:
Docker
Docker Compose

