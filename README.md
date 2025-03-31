# Fuel Price Notifier

Fuel Price Notifier is a Golang application designed to periodically update fuel prices by leveraging Go's powerful concurrency features. Developed as part of a university course, this project demonstrates how to effectively use goroutines to handle periodic tasks and maintain a clean, modular code structure.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Overview
The Fuel Price Notifier periodically retrieves updated fuel price information and processes it concurrently using goroutines. The application is designed to be scalable and serves as a practical example of building real-time data processing systems in Go.

## Features
- **Periodic Updates:** Continuously fetches fuel price data at set intervals.
- **Concurrent Processing:** Utilizes goroutines to handle multiple tasks efficiently.
- **Modular Architecture:** Organized into distinct packages such as `context`, `fluctuations`, `gas_stations`, `location`, `routers`, and `users`.
- **Ease of Extension:** A clear and maintainable codebase that can be extended to integrate additional features or data sources.

## Installation
1. **Prerequisites:** Ensure you have [Go](https://golang.org/dl/) installed on your system.
2. **Clone the repository:**

```bash
git clone https://github.com/TeoMatosevic/fuel-price-notifier.git
cd fuel-price-notifier
```

3. **Build the application:**

```bash
go build -o fuel-price-notifier
```

## Usage
- **Running with Go:**

```bash
go run main.go
```

- **Running the Built Binary:**

```bash
./fuel-price-notifier
```

The application will start fetching and processing fuel price updates based on the preconfigured intervals. Adjust any settings as needed within the code files.

## Project Structure
| Directory       | Description                                                         |    
|-----------------|---------------------------------------------------------------------|
| `context`       | Manages application context and configurations                      |
| `fluctuations`  | Contains logic for handling fuel price variations and notifications |
| `gas_stations`  | Holds models and functions related to gas station data              |
| `location`      | Provides geolocation functionalities                                |
| `routers`       | Defines HTTP routes and server setup                                |
| `users`         | Manages user-related operations and data                            |

The main entry point of the application is in `main.go`, which orchestrates the initialization and scheduling of periodic tasks.

## License
This project was developed as an academic exercise. If a license is needed, please add it here.

## Acknowledgments
- This project is inspired by real-world applications that harness concurrent programming in Go.
- Sincere thanks to peers and instructors who provided insights and feedback during development.
