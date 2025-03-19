# Traffic Monitor

A Go-based traffic monitoring system that detects abnormal traffic increases, with special support for lunar festivals.

## Features

- Monitors traffic data for specific modules and IDCs
- Detects abnormal traffic increases by comparing with historical data
- Special handling for lunar festivals (e.g., Spring Festival, Mid-Autumn Festival)
- Configurable threshold for traffic increase detection
- CSV-based data storage
- Console-based notifications (extensible to other notification channels)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/whichonezhang/traffic_monitor.git
cd traffic_monitor
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o monitor ./cmd/monitor
```

## Usage

Run the monitor with the following command:

```bash
./monitor -module=<module> -idc=<idc> [-threshold=<threshold>]
```

Parameters:
- `module`: Name of the module to monitor (required)
- `idc`: Name of the IDC to monitor (required)
- `threshold`: Threshold for traffic increase detection (default: 0.5, meaning 50% increase)

Example:
```bash
./monitor -module=api -idc=us-west -threshold=0.3
```

## Data Format

The system expects traffic data in CSV format with the following structure:

```csv
timestamp,requests
2024-02-10 00:00:00,100
2024-02-10 00:01:00,120
...
```

Data files should be placed in the `data` directory with the following naming convention:
```
data/<module>_<idc>_<YYYYMMDD>.csv
```

## Lunar Festival Support

The system automatically detects and handles the following lunar festivals:
- Spring Festival (春节)
- Lantern Festival (元宵节)
- Dragon Boat Festival (端午节)
- Mid-Autumn Festival (中秋节)

When monitoring traffic during these festivals, the system compares the current traffic with the same festival from the previous year.

## Development

### Project Structure

```
traffic_monitor/
├── cmd/
│   └── monitor/
│       └── main.go
├── internal/
│   ├── calendar/
│   │   └── lunar.go
│   ├── data/
│   │   └── provider.go
│   ├── monitor/
│   │   └── monitor.go
│   ├── notification/
│   │   └── notifier.go
│   └── types/
│       ├── notification.go
│       └── traffic.go
├── go.mod
└── README.md
```

### Adding New Features

1. **New Data Source**: Implement the `data.Provider` interface
2. **New Notification Channel**: Implement the `notification.Notifier` interface
3. **New Festival**: Add the festival date to the `festivalMap` in `calendar.LunarCalendar`

## License

MIT License 