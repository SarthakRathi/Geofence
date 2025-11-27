package geo

import "geofence/models"

func IsPointInPolygon(point models.Coordinate, polygon []models.Coordinate) bool {
	intersectCount := 0
	j := len(polygon) - 1 // The last vertex

	for i := 0; i < len(polygon); i++ {
		vertex1 := polygon[i]
		vertex2 := polygon[j]

		if ((vertex1.Lat > point.Lat) != (vertex2.Lat > point.Lat)) &&
			(point.Lon < (vertex2.Lon-vertex1.Lon)*(point.Lat-vertex1.Lat)/(vertex2.Lat-vertex1.Lat)+vertex1.Lon) {
			intersectCount++
		}
		j = i
	}

	return intersectCount%2 == 1
}
