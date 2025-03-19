package analyzer

import (
	"math"
	"time"

	"github.com/whichonezhang/traffic_monitor/internal/types"
)

// TimeSeriesAnalyzer handles advanced time series analysis
type TimeSeriesAnalyzer struct {
	data []types.TrafficData
}

// NewTimeSeriesAnalyzer creates a new time series analyzer
func NewTimeSeriesAnalyzer(data []types.TrafficData) *TimeSeriesAnalyzer {
	return &TimeSeriesAnalyzer{
		data: data,
	}
}

// Decompose decomposes the time series into trend, seasonal, and residual components
func (a *TimeSeriesAnalyzer) Decompose() (trend, seasonal, residual []float64, err error) {
	if len(a.data) < 2 {
		return nil, nil, nil, nil
	}

	// Extract values and timestamps
	values := make([]float64, len(a.data))
	timestamps := make([]time.Time, len(a.data))
	for i, d := range a.data {
		values[i] = d.Requests
		timestamps[i] = d.Timestamp
	}

	// Calculate trend using moving average
	windowSize := 7 // 7-point moving average
	trend = make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		start := max(0, i-windowSize/2)
		end := min(len(values), i+windowSize/2+1)
		sum := 0.0
		for j := start; j < end; j++ {
			sum += values[j]
		}
		trend[i] = sum / float64(end-start)
	}

	// Calculate seasonal component
	seasonal = make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		seasonal[i] = values[i] - trend[i]
	}

	// Calculate residual component
	residual = make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		residual[i] = values[i] - trend[i] - seasonal[i]
	}

	return trend, seasonal, residual, nil
}

// DetectAnomalies detects anomalies using statistical methods
func (a *TimeSeriesAnalyzer) DetectAnomalies() ([]int, error) {
	if len(a.data) < 2 {
		return nil, nil
	}

	// Extract values
	values := make([]float64, len(a.data))
	for i, d := range a.data {
		values[i] = d.Requests
	}

	// Calculate mean and standard deviation
	mean := calculateMean(values)
	stdDev := calculateStdDev(values, mean)

	// Detect anomalies (values outside 3 standard deviations)
	var anomalies []int
	for i, v := range values {
		if math.Abs(v-mean) > 3*stdDev {
			anomalies = append(anomalies, i)
		}
	}

	return anomalies, nil
}

// Forecast predicts future values using ARIMA model
func (a *TimeSeriesAnalyzer) Forecast(steps int) ([]float64, error) {
	if len(a.data) < 2 {
		return nil, nil
	}

	// Extract values
	values := make([]float64, len(a.data))
	for i, d := range a.data {
		values[i] = d.Requests
	}

	// Simple ARIMA(1,1,1) implementation
	// This is a simplified version - in production, you might want to use a proper ARIMA library
	forecast := make([]float64, steps)
	lastValue := values[len(values)-1]
	lastDiff := values[len(values)-1] - values[len(values)-2]
	phi := 0.7 // AR coefficient

	for i := 0; i < steps; i++ {
		// ARIMA(1,1,1) formula
		forecast[i] = lastValue + phi*lastDiff
		lastDiff = forecast[i] - lastValue
		lastValue = forecast[i]
	}

	return forecast, nil
}

// CalculateSeasonality calculates the seasonal pattern in the data
func (a *TimeSeriesAnalyzer) CalculateSeasonality(period int) ([]float64, error) {
	if len(a.data) < period {
		return nil, nil
	}

	// Extract values
	values := make([]float64, len(a.data))
	for i, d := range a.data {
		values[i] = d.Requests
	}

	// Calculate seasonal indices
	seasonalIndices := make([]float64, period)
	counts := make([]int, period)

	for i, v := range values {
		idx := i % period
		seasonalIndices[idx] += v
		counts[idx]++
	}

	// Calculate average seasonal pattern
	for i := 0; i < period; i++ {
		if counts[i] > 0 {
			seasonalIndices[i] /= float64(counts[i])
		}
	}

	return seasonalIndices, nil
}

// Helper functions
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) < 2 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(values)-1))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
