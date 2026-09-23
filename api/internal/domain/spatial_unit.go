package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SpatialUnit represents a land parcel with PostGIS polygon geometry (LADM LA_SpatialUnit)
// Adheres to the Fit-For-Purpose (FFP) land administration framework
type SpatialUnit struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	GeomGeoJSON    string        `json:"geom_geojson" db:"geom"`
	AreaSqm        float64       `json:"area_sqm" db:"area_sqm"`
	AccuracyMeters float64       `json:"accuracy_meters" db:"accuracy_meters"`
	CaptureMethod  CaptureMethod `json:"capture_method" db:"capture_method"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
}

// NewSpatialUnit creates and validates a new SpatialUnit parcel (Factory Method pattern)
func NewSpatialUnit(geomGeoJSON string, accuracyMeters float64, captureMethod CaptureMethod) (*SpatialUnit, error) {
	geomGeoJSON = strings.TrimSpace(geomGeoJSON)
	if geomGeoJSON == "" {
		return nil, errors.New("parcel geometry (GeoJSON) cannot be empty")
	}

	if accuracyMeters < 0 {
		return nil, errors.New("accuracy_meters cannot be negative")
	}

	if !captureMethod.IsValid() {
		return nil, errors.New("invalid capture_method")
	}

	now := time.Now().UTC()
	return &SpatialUnit{
		ID:             uuid.New(),
		GeomGeoJSON:    geomGeoJSON,
		AccuracyMeters: accuracyMeters,
		CaptureMethod:  captureMethod,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
