package types

import "time"

// TrafficData represents traffic data for a specific time period
type TrafficData struct {
	Timestamp time.Time
	Requests  float64
}
