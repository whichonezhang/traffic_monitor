package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/whichonezhang/traffic_monitor/internal/monitor"
	"go.uber.org/zap"
)

func main() {
	// Parse command line flags
	module := flag.String("module", "", "Module name to monitor")
	idc := flag.String("idc", "", "IDC name to monitor")
	threshold := flag.Float64("threshold", 0.5, "Threshold for traffic increase (0.5 = 50%)")
	flag.Parse()

	if *module == "" || *idc == "" {
		fmt.Println("Usage: monitor -module=<module> -idc=<idc> [-threshold=<threshold>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create monitor instance
	m := monitor.NewMonitor(*threshold, logger)

	// Run monitoring
	currentDate := time.Now()
	if err := m.RunMonitoring(*module, *idc, currentDate); err != nil {
		logger.Fatal("Monitoring failed", zap.Error(err))
	}

	logger.Info("Monitoring completed successfully")
}
