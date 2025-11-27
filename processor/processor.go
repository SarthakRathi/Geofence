package processor

import (
	"geofence/geo"
	"geofence/models"
	"geofence/store"
	"log"
	"time"
)

type EventProcessor struct {
	store store.VehicleStore
	zones []models.Zone
}

func NewEventProcessor(s store.VehicleStore, zones []models.Zone) *EventProcessor {
	return &EventProcessor{
		store: s,
		zones: zones,
	}
}

func (p *EventProcessor) ProcessEvent(evt models.LocationEvent) {
	point := models.Coordinate{Lat: evt.Lat, Lon: evt.Lon}

	currentZoneID := ""
	for _, zone := range p.zones {
		if geo.IsPointInPolygon(point, zone.Boundary) {
			currentZoneID = zone.ID
			break
		}
	}

	prevStatus, exists := p.store.GetVehicle(evt.VehicleID)
	prevZoneID := ""
	if exists {
		prevZoneID = prevStatus.CurrentZoneID
	}

	if prevZoneID != currentZoneID {
		p.handleTransition(evt.VehicleID, prevZoneID, currentZoneID)
	}

	p.store.UpdateVehicle(evt.VehicleID, currentZoneID, time.Unix(evt.Timestamp, 0))
}

func (p *EventProcessor) handleTransition(vehicleID, oldZone, newZone string) {
	if oldZone == "" && newZone != "" {
		log.Printf("[ALERT] Vehicle %s ENTERED zone %s", vehicleID, newZone)
	} else if oldZone != "" && newZone == "" {
		log.Printf("[ALERT] Vehicle %s EXITED zone %s", vehicleID, oldZone)
	} else if oldZone != "" && newZone != "" {
		log.Printf("[ALERT] Vehicle %s CROSSED from %s to %s", vehicleID, oldZone, newZone)
	}
}
