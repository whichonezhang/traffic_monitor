package types

import "time"

// Notification represents a traffic increase notification
type Notification struct {
	Module         string
	IDC            string
	CurrentDate    time.Time
	HistoricalDate time.Time
	Period         string
	Increase       float64
	CurrentMean    float64
	HistoricalMean float64
	Festival       string
	Anomalies      []int
	Forecast       []float64
}
