package store

import (
	"geofence/models"
	"sync"
	"time"
)

type VehicleStore interface {
	UpdateVehicle(id string, zoneID string, seenAt time.Time)
	GetVehicle(id string) (models.VehicleStatus, bool)
}

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string]models.VehicleStatus
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]models.VehicleStatus),
	}
}

func (s *InMemoryStore) UpdateVehicle(id string, zoneID string, seenAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = models.VehicleStatus{
		VehicleID:     id,
		CurrentZoneID: zoneID,
		LastSeen:      seenAt,
	}
}

func (s *InMemoryStore) GetVehicle(id string) (models.VehicleStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[id]
	return val, ok
}
