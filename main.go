package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"geofence/models"
	"geofence/processor"
	"geofence/store"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var cityCenterZone = models.Zone{
	ID:   "zone_1",
	Name: "City Center",
	Boundary: []models.Coordinate{
		{Lat: 10.0, Lon: 10.0}, {Lat: 10.0, Lon: 20.0},
		{Lat: 20.0, Lon: 20.0}, {Lat: 20.0, Lon: 10.0},
	},
}

func main() {
	go startServer()

	time.Sleep(1 * time.Second)

	runAutomatedTests()

	fmt.Println("\n[SYSTEM] Tests complete. Server is still running on :8080")
	select {}
}

func startServer() {
	memStore := store.NewInMemoryStore()
	proc := processor.NewEventProcessor(memStore, []models.Zone{cityCenterZone})

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var evt models.LocationEvent
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		proc.ProcessEvent(evt)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/vehicle/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/vehicle/")
		status, found := memStore.GetVehicle(id)
		if !found {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(status)
	})

	log.Println("[SERVER] Geofence Service Started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func runAutomatedTests() {
	fmt.Println("\n----------------------------------------")
	fmt.Println("   STARTING AUTOMATED GEOFENCE DEMO")
	fmt.Println("----------------------------------------")
	vehicleID := "taxi-demo-01"

	fmt.Println("\n[TESTER] 1. Vehicle moves to (15,15) -> EXPECT: ENTER Zone")
	sendPing(vehicleID, 15.0, 15.0)
	time.Sleep(500 * time.Millisecond)

	fmt.Println("[TESTER]    Checking API status...")
	printStatus(vehicleID)

	fmt.Println("\n[TESTER] 2. Vehicle moves to (16,16) -> EXPECT: NO ALERT (Still inside)")
	sendPing(vehicleID, 16.0, 16.0)
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n[TESTER] 3. Vehicle moves to (50,50) -> EXPECT: EXIT Zone")
	sendPing(vehicleID, 50.0, 50.0)
	time.Sleep(500 * time.Millisecond)

	fmt.Println("[TESTER]    Checking API status...")
	printStatus(vehicleID)
	fmt.Println("\n----------------------------------------")
}

// -- Helpers --

func sendPing(id string, lat, lon float64) {
	payload := map[string]interface{}{
		"vehicle_id": id, "lat": lat, "lon": lon, "timestamp": time.Now().Unix(),
	}
	data, _ := json.Marshal(payload)
	http.Post("http://localhost:8080/event", "application/json", bytes.NewBuffer(data))
}

func printStatus(id string) {
	resp, err := http.Get("http://localhost:8080/vehicle/" + id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("        > API Response: %s\n", string(body))
}
