package analyzer

import (
	"math"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/whichonezhang/traffic_monitor/internal/types"
)

func TestTimeSeriesAnalyzer(t *testing.T) {
	// Create test data with known patterns
	data := generateTestData()
	analyzer := NewTimeSeriesAnalyzer(data)

	// Test decomposition
	trend, seasonal, residual, err := analyzer.Decompose()
	assert.NoError(t, err)
	assert.NotNil(t, trend)
	assert.NotNil(t, seasonal)
	assert.NotNil(t, residual)
	assert.Equal(t, len(data), len(trend))
	assert.Equal(t, len(data), len(seasonal))
	assert.Equal(t, len(data), len(residual))

	// Test anomaly detection
	anomalies, err := analyzer.DetectAnomalies()
	assert.NoError(t, err)
	assert.NotNil(t, anomalies)
	// Should detect the spike at index 20
	assert.Contains(t, anomalies, 20)

	// Test forecasting
	forecast, err := analyzer.Forecast(5)
	assert.NoError(t, err)
	assert.NotNil(t, forecast)
	assert.Equal(t, 5, len(forecast))

	// Test seasonality calculation
	seasonalIndices, err := analyzer.CalculateSeasonality(24) // 24-hour period
	assert.NoError(t, err)
	assert.NotNil(t, seasonalIndices)
	assert.Equal(t, 24, len(seasonalIndices))
}

// generateTestData creates test data with known patterns
func generateTestData() []types.TrafficData {
	data := make([]types.TrafficData, 48) // 48 hours of data
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 48; i++ {
		// Base trend: increasing over time
		trend := 100.0 + float64(i)*2.0

		// Seasonal pattern: daily cycle
		hour := i % 24
		seasonal := 20.0 * math.Sin(2.0*math.Pi*float64(hour)/24.0)

		// Add a spike at hour 20
		spike := 0.0
		if i == 20 {
			spike = 100.0
		}

		// Add some random noise
		noise := rand.Float64()*10.0 - 5.0

		data[i] = types.TrafficData{
			Timestamp: baseTime.Add(time.Duration(i) * time.Hour),
			Requests:  trend + seasonal + spike + noise,
		}
	}

	return data
}
