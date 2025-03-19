package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/whichonezhang/traffic_monitor/internal/types"
)

// Provider interface defines the methods for data access
type Provider interface {
	GetData(module, idc string, date time.Time) ([]types.TrafficData, error)
	SaveData(module, idc string, date time.Time, data []types.TrafficData) error
}

// FileProvider implements Provider interface using CSV files
type FileProvider struct {
	dataDir string
}

// NewProvider creates a new data provider instance
func NewProvider() Provider {
	return &FileProvider{
		dataDir: "data",
	}
}

// GetData retrieves traffic data for a specific module, IDC, and date
func (p *FileProvider) GetData(module, idc string, date time.Time) ([]types.TrafficData, error) {
	filename := fmt.Sprintf("%s/%s_%s_%s.csv", p.dataDir, module, idc, date.Format("20060102"))

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV data: %w", err)
	}

	if len(records) < 1 { // Need at least one data row
		return nil, fmt.Errorf("invalid file format: insufficient data")
	}

	// Verify record format
	if len(records[0]) != 1443 { // module, idc, date, and 1440 minutes of data
		return nil, fmt.Errorf("invalid file format: expected 1443 columns, got %d", len(records[0]))
	}

	// Verify module and idc match
	if records[0][0] != module || records[0][1] != idc {
		return nil, fmt.Errorf("module/idc mismatch: expected %s/%s, got %s/%s", module, idc, records[0][0], records[0][1])
	}

	// Parse date from record
	headerDate, err := time.Parse("20060102", records[0][2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse date from record: %w", err)
	}

	// Verify date matches (compare only year, month, day)
	if headerDate.Year() != date.Year() || headerDate.Month() != date.Month() || headerDate.Day() != date.Day() {
		return nil, fmt.Errorf("date mismatch: expected %s, got %s", date.Format("20060102"), headerDate.Format("20060102"))
	}

	// Process data
	var data []types.TrafficData
	baseTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

	// Skip module, idc, and date columns
	for j := 3; j < len(records[0]); j++ {
		requests, err := strconv.ParseFloat(records[0][j], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse requests at column %d: %w", j+1, err)
		}

		timestamp := baseTime.Add(time.Duration(j-3) * time.Minute)
		data = append(data, types.TrafficData{
			Timestamp: timestamp,
			Requests:  requests,
		})
	}

	return data, nil
}

// SaveData saves traffic data to a CSV file
func (p *FileProvider) SaveData(module, idc string, date time.Time, data []types.TrafficData) error {
	filename := fmt.Sprintf("%s/%s_%s_%s.csv", p.dataDir, module, idc, date.Format("20060102"))

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create data file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Create data row
	dataRow := make([]string, 1443)
	dataRow[0] = module
	dataRow[1] = idc
	dataRow[2] = date.Format("20060102")

	// Group data by minute
	minuteData := make(map[int]float64)
	for _, d := range data {
		minute := d.Timestamp.Hour()*60 + d.Timestamp.Minute()
		minuteData[minute] = d.Requests
	}

	// Fill in the 1440 minutes of data
	for i := 0; i < 1440; i++ {
		if requests, exists := minuteData[i]; exists {
			dataRow[i+3] = fmt.Sprintf("%.2f", requests)
		} else {
			dataRow[i+3] = "0.00" // Fill missing data with 0
		}
	}

	// Write row
	if err := writer.Write(dataRow); err != nil {
		return fmt.Errorf("failed to write data row: %w", err)
	}

	return nil
}
