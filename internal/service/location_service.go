package service

import (
	"context"
	"errors"
	"log"
	"transjakarta-fleet/internal/domain"
	"transjakarta-fleet/internal/rabbitmq"
	"transjakarta-fleet/internal/repository"
)

type LocationService struct {
	repo      repository.LocationRepository
	rabbitPub *rabbitmq.Publisher
}

func NewLocationService(repo repository.LocationRepository, rabbitPub *rabbitmq.Publisher) *LocationService {
	return &LocationService{
		repo:      repo,
		rabbitPub: rabbitPub,
	}
}

func (s *LocationService) ProcessLocation(ctx context.Context, loc domain.LocationMessage) error {
	if loc.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}
	if loc.Latitude < -90 || loc.Latitude > 90 {
		return errors.New("invalid latitude")
	}
	if loc.Longitude < -180 || loc.Longitude > 180 {
		return errors.New("invalid longitude")
	}
	if loc.Timestamp <= 0 {
		return errors.New("invalid timestamp")
	}

	if err := s.repo.Insert(ctx, loc); err != nil {
		return err
	}

	if IsInsideGeofence(loc.Latitude, loc.Longitude) {
		var event domain.GeofenceEvent
		event.VehicleID = loc.VehicleID
		event.Event = "geofence_entry"
		event.Location.Latitude = loc.Latitude
		event.Location.Longitude = loc.Longitude
		event.Timestamp = loc.Timestamp

		if s.rabbitPub != nil {
			if err := s.rabbitPub.Publish(event); err != nil {
				log.Printf("failed to publish geofence event: %v", err)
			} else {
				log.Printf("geofence event published for vehicle=%s", loc.VehicleID)
			}
		}
	}

	return nil
}

func (s *LocationService) GetLatest(ctx context.Context, vehicleID string) (*domain.LocationMessage, error) {
	return s.repo.GetLatest(ctx, vehicleID)
}

func (s *LocationService) GetHistory(ctx context.Context, vehicleID string, start, end int64) ([]domain.LocationMessage, error) {
	return s.repo.GetHistory(ctx, vehicleID, start, end)
}
