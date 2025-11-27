# Geofence Event Processing Service

A high-performance, concurrent Go service designed to track vehicles and detect geofence entry/exit events in real-time. Built as a coding challenge solution within a 2-hour constraint.

## 🚀 Architecture Overview

The system follows a **Modular Monolith** architecture using the **Standard Go Project Layout**.

* **Ingestion Layer (`cmd/server`):** HTTP handlers that accept high-throughput GPS telemetry.
* **Business Logic (`internal/processor`):** Determines state transitions (Enter/Exit) by comparing the new location against the vehicle's last known state.
* **Geometry Engine (`internal/geo`):** Implements the Ray Casting algorithm to perform Point-in-Polygon checks.
* **Storage (`internal/store`):** A thread-safe in-memory repository using `sync.RWMutex` to handle concurrent reads/writes.

---

## 🛠️ Design Decisions & Trade-offs

Given the **2-hour time limit**, specific engineering tradeoffs were made to balance speed of implementation with code quality:

| Decision | Choice | Rationale (Engineering Judgment) | Production Alternative |
|:---------|:-------|:--------------------------------|:-----------------------|
| **Storage** | In-Memory Map | Fastest to implement; 0ms latency; no external dependencies (Docker/Postgres). | **Redis** (for speed) or **PostGIS** (for complex spatial queries). |
| **Concurrency** | `sync.RWMutex` | Go Maps are not thread-safe. A Mutex ensures data integrity during concurrent HTTP requests. | Sharding data by Vehicle ID to reduce lock contention. |
| **Geometry** | Ray Casting | Simple, standard algorithm for "Point in Polygon". No heavy GIS libraries needed. | **S2 Geometry** or **H3** (Uber's index) for fast lookups at scale. |
| **Search** | Linear Scan O(N) | Acceptable for <1,000 zones. | **R-Tree** or **QuadTree** index to query zones in O(log N). |

---

## 🏃 How to Run

### Prerequisites
* Go 1.18+
* (Optional) VS Code with REST Client extension

### 1. Start the Server

The server initializes with a default zone: **"City Center"** (Square from Lat/Lon 10.0 to 20.0).

```bash
go run cmd/server/main.go
```

You should see:

```
[SERVER] Geofence Service Started on :8080
[SERVER] Loaded 1 Zone: City Center (ID: zone_1)
```

### 2. Run Tests

You can use the included `requests.http` file if you use VS Code, or use the curl commands below.

---

## 🔌 API Documentation

### 1. Send Location Event

Ingests a GPS ping. If the vehicle changes zones (enters/exits), a log is emitted to stdout.

**Endpoint:** `POST /event`

**Content-Type:** `application/json`

**Request Payload:**

```json
{
  "vehicle_id": "taxi-01",
  "lat": 15.0,
  "lon": 15.0,
  "timestamp": 1672531200
}
```

**Response:** `200 OK`

**Server Log Output (Example):**

```
[ALERT] Vehicle taxi-01 ENTERED zone zone_1
```

### 2. Get Vehicle Status

Retrieves the current known location and zone status of a vehicle.

**Endpoint:** `GET /vehicle/{vehicle_id}`

**Response Payload (Inside Zone):**

```json
{
  "vehicle_id": "taxi-01",
  "current_zone_id": "zone_1",
  "last_seen": "2023-01-01T00:00:00Z"
}
```

**Response Payload (Outside Zone):**

```json
{
  "vehicle_id": "taxi-01",
  "current_zone_id": "", 
  "last_seen": "2023-01-01T00:00:00Z"
}
```

---

## 🧪 Testing Scenarios

We assume a square zone defined by (10,10) to (20,20).

1. **Enter Zone:** Send point (15, 15).
   * Result: Log `ENTERED` zone. Status returns `zone_1`.

2. **Move Inside:** Send point (16, 16).
   * Result: No log (state unchanged). Status returns `zone_1`.

3. **Exit Zone:** Send point (50, 50).
   * Result: Log `EXITED` zone. Status returns empty `""` zone.

4. **Edge Case:** Send point (10, 10) (On the boundary).
   * Result: Ray casting usually treats bottom/left edges as inclusive. Treated as "Inside".

---

## 🔮 Future Improvements (Scaling Strategy)

If we had more time or needed to scale to 100k+ vehicles:

* **Persistence:** Replace in-memory map with Redis to survive server restarts.

* **Event Bus:** Instead of logging to stdout, publish events to Apache Kafka or AWS SQS so downstream services (Push Notifications, Billing) can consume them.

* **Spatial Indexing:** A linear loop through zones is slow if we have 10,000 zones. I would implement an R-Tree or QuadTree to find relevant zones efficiently.

* **Protobuf/gRPC:** Switch from JSON to Protobuf for smaller payload size and faster serialization on mobile networks.

---

## 📝 License

This project was created as a coding challenge solution.

## 🤝 Contributing

This is a demonstration project. Feel free to fork and experiment with improvements!
