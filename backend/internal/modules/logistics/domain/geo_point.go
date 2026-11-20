package domain

// GeoPoint represents WGS-84 spatial coordinates with an optional depot or site label.
type GeoPoint struct {
	Latitude     float64
	Longitude    float64
	LocationName string
}
