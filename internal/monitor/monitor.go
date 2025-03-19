package monitor

import (
	"fmt"
	"time"

	"github.com/whichonezhang/traffic_monitor/internal/analyzer"
	"github.com/whichonezhang/traffic_monitor/internal/calendar"
	"github.com/whichonezhang/traffic_monitor/internal/data"
	"github.com/whichonezhang/traffic_monitor/internal/notification"
	"github.com/whichonezhang/traffic_monitor/internal/types"
	"go.uber.org/zap"
)

// TrafficData represents the traffic data for a specific time period
type TrafficData struct {
	Timestamp time.Time
	Requests  float64
}

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

// Monitor handles traffic monitoring and anomaly detection
type Monitor struct {
	threshold    float64
	logger       *zap.Logger
	calendar     *calendar.LunarCalendar
	dataProvider data.Provider
	notifier     notification.Notifier
}

// NewMonitor creates a new traffic monitor instance
func NewMonitor(threshold float64, logger *zap.Logger) *Monitor {
	return &Monitor{
		threshold:    threshold,
		logger:       logger,
		calendar:     calendar.NewLunarCalendar(),
		dataProvider: data.NewProvider(),
		notifier:     notification.NewNotifier(),
	}
}

// IsLunarFestival checks if a given date is a lunar festival
func (m *Monitor) IsLunarFestival(date time.Time) (string, bool) {
	return m.calendar.GetFestival(date)
}

// GetPreviousLunarFestivalDate gets the date of the same lunar festival from the previous year
func (m *Monitor) GetPreviousLunarFestivalDate(currentDate time.Time, festival string) (time.Time, error) {
	return m.calendar.GetPreviousFestivalDate(currentDate, festival)
}

// CalculateIncrease calculates the percentage increase between current and previous values
func (m *Monitor) CalculateIncrease(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return (current - previous) / previous
}

// CompareTraffic compares current traffic with historical traffic
func (m *Monitor) CompareTraffic(current, historical []types.TrafficData, timeDiff string) (float64, bool) {
	currentMean := calculateMean(current)
	historicalMean := calculateMean(historical)

	increase := m.CalculateIncrease(currentMean, historicalMean)

	if increase > m.threshold {
		m.logger.Warn("Significant traffic increase detected",
			zap.String("period", timeDiff),
			zap.Float64("increase", increase))
		return increase, true
	}
	return 0, false
}

// MonitorTraffic monitors traffic changes for different time periods
func (m *Monitor) MonitorTraffic(module, idc string, currentDate time.Time) ([]types.Notification, error) {
	var notifications []types.Notification

	// Get current data
	currentData, err := m.dataProvider.GetData(module, idc, currentDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get current data: %w", err)
	}

	// Create time series analyzer for current data
	currentAnalyzer := analyzer.NewTimeSeriesAnalyzer(currentData)

	// Detect anomalies in current data
	anomalies, err := currentAnalyzer.DetectAnomalies()
	if err != nil {
		return nil, fmt.Errorf("failed to detect anomalies: %w", err)
	}

	// Forecast future traffic
	forecast, err := currentAnalyzer.Forecast(24) // Forecast next 24 hours
	if err != nil {
		return nil, fmt.Errorf("failed to forecast traffic: %w", err)
	}

	// Check if current date is a lunar festival
	if festival, isFestival := m.IsLunarFestival(currentDate); isFestival {
		previousDate, err := m.GetPreviousLunarFestivalDate(currentDate, festival)
		if err != nil {
			return nil, fmt.Errorf("failed to get previous festival date: %w", err)
		}

		historicalData, err := m.dataProvider.GetData(module, idc, previousDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get historical data: %w", err)
		}

		if increase, significant := m.CompareTraffic(currentData, historicalData, fmt.Sprintf("Previous %s", festival)); significant {
			notifications = append(notifications, types.Notification{
				Module:         module,
				IDC:            idc,
				CurrentDate:    currentDate,
				HistoricalDate: previousDate,
				Period:         fmt.Sprintf("Previous %s", festival),
				Increase:       increase,
				CurrentMean:    calculateMean(currentData),
				HistoricalMean: calculateMean(historicalData),
				Festival:       festival,
				Anomalies:      anomalies,
				Forecast:       forecast,
			})
		}
	}

	// Compare with regular time periods
	timePeriods := []struct {
		name     string
		duration time.Duration
	}{
		{"1 day ago", 24 * time.Hour},
		{"7 days ago", 7 * 24 * time.Hour},
		{"30 days ago", 30 * 24 * time.Hour},
		{"1 year ago", 365 * 24 * time.Hour},
	}

	for _, period := range timePeriods {
		historicalDate := currentDate.Add(-period.duration)
		historicalData, err := m.dataProvider.GetData(module, idc, historicalDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get historical data: %w", err)
		}

		if increase, significant := m.CompareTraffic(currentData, historicalData, period.name); significant {
			notifications = append(notifications, types.Notification{
				Module:         module,
				IDC:            idc,
				CurrentDate:    currentDate,
				HistoricalDate: historicalDate,
				Period:         period.name,
				Increase:       increase,
				CurrentMean:    calculateMean(currentData),
				HistoricalMean: calculateMean(historicalData),
				Anomalies:      anomalies,
				Forecast:       forecast,
			})
		}
	}

	return notifications, nil
}

// RunMonitoring runs the traffic monitoring process
func (m *Monitor) RunMonitoring(module, idc string, currentDate time.Time) error {
	notifications, err := m.MonitorTraffic(module, idc, currentDate)
	if err != nil {
		return fmt.Errorf("monitoring failed: %w", err)
	}

	for _, notification := range notifications {
		if err := m.notifier.Send(notification); err != nil {
			m.logger.Error("Failed to send notification",
				zap.Error(err),
				zap.Any("notification", notification))
		}
	}

	return nil
}

// Helper function to calculate mean of traffic data
func calculateMean(data []types.TrafficData) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, d := range data {
		sum += d.Requests
	}
	return sum / float64(len(data))
}
