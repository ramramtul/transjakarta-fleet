package repository

import (
	"context"
	"database/sql"
	"transjakarta-fleet/internal/domain"
)

type LocationRepository interface {
	Insert(ctx context.Context, loc domain.LocationMessage) error
	GetLatest(ctx context.Context, vehicleID string) (*domain.LocationMessage, error)
	GetHistory(ctx context.Context, vehicleID string, start, end int64) ([]domain.LocationMessage, error)
}

type locationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Insert(ctx context.Context, loc domain.LocationMessage) error {
	query := `
		INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp)
	return err
}

func (r *locationRepository) GetLatest(ctx context.Context, vehicleID string) (*domain.LocationMessage, error) {
	query := `
		SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var loc domain.LocationMessage
	err := r.db.QueryRowContext(ctx, query, vehicleID).
		Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp)
	if err != nil {
		return nil, err
	}

	return &loc, nil
}

func (r *locationRepository) GetHistory(ctx context.Context, vehicleID string, start, end int64) ([]domain.LocationMessage, error) {
	query := `
		SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1
		  AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, vehicleID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.LocationMessage
	for rows.Next() {
		var loc domain.LocationMessage
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return nil, err
		}
		result = append(result, loc)
	}

	return result, nil
}
