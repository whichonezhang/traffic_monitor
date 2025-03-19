package monitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/whichonezhang/traffic_monitor/internal/types"
	"go.uber.org/zap"
)

func TestMonitorTraffic(t *testing.T) {
	// Create test logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create monitor instance
	m := NewMonitor(0.5, logger)

	// Test case 1: Normal traffic (no significant increase)
	currentDate := time.Date(2024, 2, 10, 0, 0, 0, 0, time.Local)
	notifications, err := m.MonitorTraffic("api", "us-west", currentDate)
	assert.NoError(t, err)
	assert.Empty(t, notifications, "Should not detect any significant increase for normal traffic")

	// Test case 2: Festival traffic (should compare with previous year)
	currentDate = time.Date(2024, 2, 10, 0, 0, 0, 0, time.Local)
	notifications, err = m.MonitorTraffic("api", "us-west", currentDate)
	assert.NoError(t, err)
	assert.Empty(t, notifications, "Should not detect any significant increase for festival traffic")

	// Test case 3: High threshold (should not detect increase)
	m = NewMonitor(1.0, logger)
	notifications, err = m.MonitorTraffic("api", "us-west", currentDate)
	assert.NoError(t, err)
	assert.Empty(t, notifications, "Should not detect increase with high threshold")

	// Test case 4: Low threshold (should detect increase)
	m = NewMonitor(0.1, logger)
	notifications, err = m.MonitorTraffic("api", "us-west", currentDate)
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications, "Should detect increase with low threshold")
}

func TestCalculateIncrease(t *testing.T) {
	m := NewMonitor(0.5, nil)

	tests := []struct {
		name     string
		current  float64
		previous float64
		expected float64
	}{
		{"Normal increase", 150, 100, 0.5},
		{"No increase", 100, 100, 0.0},
		{"Decrease", 50, 100, -0.5},
		{"Zero previous", 100, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.CalculateIncrease(tt.current, tt.previous)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompareTraffic(t *testing.T) {
	m := NewMonitor(0.5, nil)

	tests := []struct {
		name           string
		current        []types.TrafficData
		historical     []types.TrafficData
		timeDiff       string
		expectIncrease bool
	}{
		{
			name: "Significant increase",
			current: []types.TrafficData{
				{Requests: 150},
				{Requests: 160},
				{Requests: 170},
			},
			historical: []types.TrafficData{
				{Requests: 100},
				{Requests: 110},
				{Requests: 120},
			},
			timeDiff:       "1 day ago",
			expectIncrease: true,
		},
		{
			name: "No significant increase",
			current: []types.TrafficData{
				{Requests: 110},
				{Requests: 120},
				{Requests: 130},
			},
			historical: []types.TrafficData{
				{Requests: 100},
				{Requests: 110},
				{Requests: 120},
			},
			timeDiff:       "1 day ago",
			expectIncrease: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			increase, significant := m.CompareTraffic(tt.current, tt.historical, tt.timeDiff)
			assert.Equal(t, tt.expectIncrease, significant)
			if significant {
				assert.Greater(t, increase, 0.5)
			} else {
				assert.LessOrEqual(t, increase, 0.5)
			}
		})
	}
}
