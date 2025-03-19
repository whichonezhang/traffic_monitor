package notification

import (
	"fmt"
	"strings"

	"github.com/whichonezhang/traffic_monitor/internal/types"
)

// Notifier interface defines the methods for sending notifications
type Notifier interface {
	Send(notification types.Notification) error
}

// ConsoleNotifier implements Notifier interface using console output
type ConsoleNotifier struct{}

// NewNotifier creates a new notifier instance
func NewNotifier() Notifier {
	return &ConsoleNotifier{}
}

// Send sends a notification about traffic changes
func (n *ConsoleNotifier) Send(notification types.Notification) error {
	message := formatNotification(notification)
	fmt.Println(message)
	return nil
}

// formatNotification formats the notification message
func formatNotification(n types.Notification) string {
	var message strings.Builder

	// Write header
	if n.Festival != "" {
		message.WriteString(fmt.Sprintf(
			"Traffic Alert for %s Festival\n",
			n.Festival))
	} else {
		message.WriteString("Traffic Alert\n")
	}

	// Write basic information
	message.WriteString(fmt.Sprintf(
		"Module: %s\n"+
			"IDC: %s\n"+
			"Current Date: %s\n"+
			"Historical Date: %s\n"+
			"Period: %s\n"+
			"Increase: %.2f%%\n"+
			"Current Mean: %.2f\n"+
			"Historical Mean: %.2f\n",
		n.Module,
		n.IDC,
		n.CurrentDate.Format("2006-01-02"),
		n.HistoricalDate.Format("2006-01-02"),
		n.Period,
		n.Increase*100,
		n.CurrentMean,
		n.HistoricalMean))

	// Write anomaly information
	if len(n.Anomalies) > 0 {
		message.WriteString("\nDetected Anomalies:\n")
		for _, idx := range n.Anomalies {
			message.WriteString(fmt.Sprintf("- Anomaly at index %d\n", idx))
		}
	}

	// Write forecast information
	if len(n.Forecast) > 0 {
		message.WriteString("\nTraffic Forecast (next 24 hours):\n")
		for i, value := range n.Forecast {
			message.WriteString(fmt.Sprintf("- Hour %d: %.2f\n", i+1, value))
		}
	}

	return message.String()
}
