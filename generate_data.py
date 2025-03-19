#!/usr/bin/env python3

import os
import random
from datetime import datetime, timedelta

def generate_data(module, idc, date):
    # Create data directory if it doesn't exist
    os.makedirs("data", exist_ok=True)

    # Generate file path
    file_path = f"data/{module}_{idc}_{date.strftime('%Y%m%d')}.csv"

    # Generate data
    data = [module, idc, date.strftime("%Y%m%d")]
    base_value = 100
    for i in range(1440):  # 24 hours * 60 minutes
        value = base_value + random.randint(-10, 10)
        data.append(f"{value:.2f}")

    # Write to file
    with open(file_path, "w") as f:
        f.write(",".join(data))

def main():
    module = "api"
    idc = "us-west"
    
    # Generate current data
    current_date = datetime(2025, 3, 19)
    generate_data(module, idc, current_date)
    
    # Generate historical data
    historical_dates = [
        current_date - timedelta(days=1),  # 1 day ago
        current_date - timedelta(days=7),  # 7 days ago
        current_date - timedelta(days=30),  # 30 days ago
        current_date - timedelta(days=365),  # 1 year ago
    ]
    
    for date in historical_dates:
        generate_data(module, idc, date)

if __name__ == "__main__":
    main() 