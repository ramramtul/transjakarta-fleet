# Transjakarta Fleet Management Backend

Backend service for a fleet management system built with **Golang**, **MQTT**, **PostgreSQL**, **RabbitMQ**, and **Docker Compose**.

This project was developed as a technical test submission for **Backend Engineer - Transjakarta**.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [System Architecture](#system-architecture)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [Prerequisites](#prerequisites)
- [How to Run](#how-to-run)
- [Available Services](#available-services)
- [API Endpoints](#api-endpoints)
- [MQTT Integration](#mqtt-integration)
- [RabbitMQ Integration](#rabbitmq-integration)
- [Database Schema](#database-schema)
- [Example Testing Flow](#example-testing-flow)
- [Notes / Design Decisions](#notes--design-decisions)
- [Author](#author)

---

## Overview

This system simulates a **vehicle fleet tracking backend** where vehicle GPS locations are sent through **MQTT**, processed by a **Golang backend**, stored in **PostgreSQL**, and monitored for **geofence events** that are published to **RabbitMQ**.

The backend also provides REST APIs to retrieve:

- the **latest location** of a vehicle
- the **history of vehicle locations** within a specific time range

All services run using **Docker Compose**.

---

## Features

- Receive vehicle location data via **MQTT**
- Validate incoming location payloads
- Store location data in **PostgreSQL**
- Provide REST APIs for:
  - latest vehicle location
  - vehicle location history
- Detect **geofence entry**
- Publish geofence events to **RabbitMQ**
- Consume geofence events via worker
- Run the full stack with **Docker Compose**
- Include a mock publisher that sends location data every **2 seconds**

---

## Tech Stack

- **Language:** Golang
- **Web Framework:** Gin
- **Database:** PostgreSQL
- **Message Broker:** RabbitMQ
- **MQTT Broker:** Eclipse Mosquitto
- **Containerization:** Docker & Docker Compose

---

## System Architecture

### High-Level Flow

```bash
Mock Publisher
   ↓
MQTT Broker (Mosquitto)
   ↓
Backend Subscriber (Golang)
   ↓
PostgreSQL
   ↓
Geofence Check
   ↓
RabbitMQ Exchange (fleet.events)
   ↓
Queue (geofence_alerts)
   ↓
Worker Consumer
```

## Project Structure
```bash
transjakarta-fleet/
├── cmd/
│   ├── app/                  # Main backend application
│   └── publisher/            # Mock MQTT publisher
│
├── internal/
│   ├── config/               # Retry helper / startup config
│   ├── database/             # PostgreSQL connection
│   ├── domain/               # Request / response / event models
│   ├── handler/              # HTTP handlers
│   ├── mqtt/                 # MQTT subscriber
│   ├── rabbitmq/             # RabbitMQ publisher + worker
│   ├── repository/           # Database access layer
│   ├── router/               # Gin router setup
│   └── service/              # Business logic
│
├── migrations/               # SQL migration files
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## How It Works
1. Mock Publisher Sends Location via MQTT
   A mock publisher service sends GPS data every 2 seconds to the topic:
   `/fleet/vehicle/{vehicle_id}/location`

    Example payload:
    ```bash
    {
      "vehicle_id": "B1234XYZ",
      "latitude": -6.2088,
      "longitude": 106.8456,
      "timestamp": 1715003456
    }
    ```
2. Backend Subscribes and Validates
   The backend subscribes to:
   `/fleet/vehicle/+/location`

   It validates:
    - vehicle_id must not be empty
    - latitude must be between -90 and 90
    - longitude must be between -180 and 180
    - timestamp must be valid

3. Data is Stored in PostgreSQL
   Validated location data is inserted into PostgreSQL table vehicle_locations

4. Backend Checks Geofence
   Each incoming location is checked against a geofence area.

   Current geofence configuration:
   ```bash
   Center Latitude: -6.2088
   Center Longitude: 106.8456
   Radius: 50 meters
   ```
   If a vehicle enters the geofence, the backend publishes an event to RabbitMQ.

5. Geofence Event is Published to RabbitMQ
   Event is published to:
   ```bash 
    Exchange: fleet.events
    Queue: geofence_alerts
    
    Example event payload:
    {
      "vehicle_id": "B1234XYZ",
      "event": "geofence_entry",
      "location": {
        "latitude": -6.2088,
        "longitude": 106.8456
      },
      "timestamp": 1715003456
    }
    ```

6. Worker Consumes Geofence Event
A worker inside the backend consumes messages from the queue and logs them.

7. REST API Exposes Data
   The backend provides APIs to retrieve:
    - latest vehicle location
    - location history by time range

## Prerequisites
Make sure you have the following installed:
- Docker
- Docker Compose

No need to install PostgreSQL, RabbitMQ, or MQTT broker manually because everything runs in Docker.

## How to Run
1. Clone Repository
    ```bash
    git clone <your-github-repository-url>
    cd transjakarta-fleet
    ```

2. Run All Services
   `docker compose up --build`

   This command will start:
    - Backend API
    - PostgreSQL
    - RabbitMQ
    - Mosquitto MQTT Broker
    - Mock MQTT Publisher

3. Wait Until Services Are Ready
   You should see logs similar to:
    ```bash
    PostgreSQL connected successfully
    RabbitMQ exchange and queue declared successfully
    RabbitMQ worker started
    Connected to MQTT broker: tcp://mosquitto:1883
    MQTT subscriber started
    HTTP server running on :8080
    ```

    You should also see publisher logs sending location data every 2 seconds.

## Available Services
| Service                | URL / Port               |
| ---------------------- | ------------------------ |
| Backend API            | `http://localhost:8081`  |
| PostgreSQL             | `localhost:5432`         |
| RabbitMQ Management UI | `http://localhost:15672` |
| MQTT Broker            | `localhost:1883`         |

## RabbitMQ Login
- Username: guest
- Password: guest

## API Endpoints
1. Health Check
    ```bash
    Request
    GET /health
    
    Example
    curl http://localhost:8081/health
    
    Response
    {
      "status": "ok"
    }
    ```

2. Get Latest Vehicle Location
    ```bash
    Request
    GET /vehicles/:vehicle_id/location
    
    Example
    curl http://localhost:8081/vehicles/B1234XYZ/location
    
    Example Response
    {
      "vehicle_id": "B1234XYZ",
      "latitude": -6.2088,
      "longitude": 106.8456,
      "timestamp": 1715003456
    }
    ```

3. Get Vehicle Location History
    ```bash
    Request
    GET /vehicles/:vehicle_id/history?start=1715000000&end=1715009999
    
    Example
    curl "http://localhost:8081/vehicles/B1234XYZ/history?start=1715000000&end=1999999999"
    
    Example Response
    [
      {
        "vehicle_id": "B1234XYZ",
        "latitude": -6.2088,
        "longitude": 106.8456,
        "timestamp": 1715003456
      },
      {
        "vehicle_id": "B1234XYZ",
        "latitude": -6.2087,
        "longitude": 106.8455,
        "timestamp": 1715003458
      }
    ]
    ```

## MQTT Integration
```bash
Topic
/fleet/vehicle/{vehicle_id}/location

Example Topic
/fleet/vehicle/B1234XYZ/location

Payload Format
{
  "vehicle_id": "B1234XYZ",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "timestamp": 1715003456
}
```

Publisher Behavior
- Sends mock GPS location every 2 seconds
- Uses slight random coordinate variation to simulate movement

## RabbitMQ Integration
```bash
Exchange
fleet.events

Queue
geofence_alerts

Routing Key
geofence.entry

Event Payload
{
  "vehicle_id": "B1234XYZ",
  "event": "geofence_entry",
  "location": {
    "latitude": -6.2088,
    "longitude": 106.8456
  },
  "timestamp": 1715003456
}
```

Worker Behavior
The backend includes a worker that consumes messages from geofence_alerts and logs them.

## Database Schema
Table: vehicle_locations
| Column     | Type             |
| ---------- | ---------------- |
| vehicle_id | TEXT             |
| latitude   | DOUBLE PRECISION |
| longitude  | DOUBLE PRECISION |
| timestamp  | BIGINT           |

Example Query
To verify inserted data manually:
`docker exec -it fleet_postgres psql -U postgres -d fleetdb`

Then run:
```bash
SELECT vehicle_id, latitude, longitude, timestamp
FROM vehicle_locations
ORDER BY timestamp DESC
LIMIT 10;
```

## Example Testing Flow
1. Start the system
   `docker compose up --build`

2. Wait a few seconds for mock data to be published
   The publisher sends data every 2 seconds.

4. Test latest location API
   `curl http://localhost:8081/vehicles/B1234XYZ/location`

5. Test location history API
   `curl "http://localhost:8081/vehicles/B1234XYZ/history?start=1715000000&end=1999999999"`
   
6. Check RabbitMQ geofence event
   Watch backend logs:
    ```bash
    geofence event published for vehicle=B1234XYZ
    RabbitMQ worker received geofence event: ...
    ```

   Or open RabbitMQ UI: `http://localhost:15672`

## Notes / Design Decisions
1. Dockerized Full Stack
   All services are containerized for easier setup and reproducibility.

2. Retry Logic on Startup
   The backend includes retry logic for:
    - PostgreSQL
    - RabbitMQ
    - MQTT

   This is useful because in Docker Compose, a service may start before its dependency is fully ready.

3. Layered Project Structure
   The code is organized into:
    - handler layer
    - service layer
    - repository layer

   This separation makes the project easier to maintain and extend.

4. Simple Geofence Detection
   Geofence detection uses the Haversine formula to calculate distance between coordinates.

5. Mock Data Publisher
   A dedicated publisher service is included to simulate real-time vehicle movement and make the system easy to test.

   This repository includes:
  
   Source code
    - Docker Compose setup
    - Dockerfile
    - SQL migration
    - README
    - Postman Collection
    - Demo-ready architecture flow

## Author

Rama Rahmatullah
Backend Engineer Technical Test Submission - Transjakarta
