package data

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/whichonezhang/traffic_monitor/internal/types"
)

func TestFileProvider(t *testing.T) {
	// Create test data
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	testModule := "api"
	testIDC := "us-west"

	// Create test data with 1440 minutes
	testData := make([]types.TrafficData, 1440)
	for i := 0; i < 1440; i++ {
		testData[i] = types.TrafficData{
			Timestamp: testDate.Add(time.Duration(i) * time.Minute),
			Requests:  float64(100 + i),
		}
	}

	// Create provider
	provider := NewProvider()

	// Test SaveData
	err := provider.SaveData(testModule, testIDC, testDate, testData)
	assert.NoError(t, err)

	// Test GetData
	retrievedData, err := provider.GetData(testModule, testIDC, testDate)
	assert.NoError(t, err)
	assert.Equal(t, len(testData), len(retrievedData))

	// Verify data matches
	for i := 0; i < len(testData); i++ {
		assert.Equal(t, testData[i].Timestamp, retrievedData[i].Timestamp)
		assert.Equal(t, testData[i].Requests, retrievedData[i].Requests)
	}

	// Clean up test file
	filename := "data/api_us-west_20240101.csv"
	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestFileProviderErrors(t *testing.T) {
	provider := NewProvider()
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)

	// Test non-existent file
	_, err := provider.GetData("nonexistent", "nonexistent", testDate)
	assert.Error(t, err)

	// Test invalid data format
	err = os.WriteFile("data/invalid.csv", []byte("invalid,data\n1,2"), 0644)
	assert.NoError(t, err)
	defer os.Remove("data/invalid.csv")

	_, err = provider.GetData("invalid", "invalid", testDate)
	assert.Error(t, err)
}
