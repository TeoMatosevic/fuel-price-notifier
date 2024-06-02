package location

import (
    "math"
)

type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

func Create(lat, lon float64) *Location {
    return &Location{lat, lon}
}

func (l *Location) Verify() bool {
    return l.Latitude >= -90 && l.Latitude <= 90 && l.Longitude >= -180 && l.Longitude <= 180
}

func (l *Location) DistanceTo(l2 *Location) float64 {
    // Haversine formula
    const R = 6371
    lat1 := l.Latitude * (math.Pi / 180)
    lat2 := l2.Latitude * (math.Pi / 180)
    d_lat := (l2.Latitude - l.Latitude) * (math.Pi / 180)
    d_lon := (l2.Longitude - l.Longitude) * (math.Pi / 180)

    a := math.Pow(math.Sin(d_lat/2), 2) + math.Cos(lat1) * math.Cos(lat2) * math.Pow(math.Sin(d_lon/2), 2)
    c := 2 * math.Asin(math.Sqrt(a))

    return R * c
}
