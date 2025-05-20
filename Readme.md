# Radar240 - ASTERIX CAT 240 Parser

This project implements a parser and decoder for ASTERIX (All Purpose Structured EUROCONTROL Surveillance Information Exchange) Category 240 data, specifically designed for radar video data processing.

## Features

- Parses ASTERIX CAT 240 messages from PCAP files
- Decompresses radar data when compression is used
- Converts polar coordinates to geographical coordinates
- Real-time data streaming via WebSocket
- Supports different video data resolutions (1, 2, 4, 8 bits per cell)

## Architecture

The project consists of several key components:

1. **Parser** ([parser/cat240](parser/cat240/Parser.go))
   - Decodes ASTERIX messages
   - Validates data structures
   - Handles different data formats

2. **Data Processing** ([utils](utils/))
   - PCAP file reading
   - Data parsing and transformation
   - Coordinate system conversion

3. **WebSocket Server** ([sender](sender/sendData240.go))
   - Real-time data streaming
   - Client connection management

## Installation

```sh
# Clone the repository
git clone <repository-url>

# Install dependencies
go mod download

# Build the project
go build
```

## Usage

### Running with Docker

```sh
docker-compose up --build
```

### Running Locally

```sh
go run main.go
```

The server will start on port 8080 and begin processing the PCAP file located in `data/ASTERIX_CAT240_1_20230517184234.pcap`.

## Data Structure

The main data structures are defined in [cat240Struct.go](parser/cat240/cat240Struct.go):

- `VideoDataItem`: Contains raw message data
- `ValidData`: Contains validated and processed data
- `BlockData`: Contains decoded geographical coordinates and intensity values

## API

### WebSocket Endpoint

```
ws://localhost:8080/radar240
```

The WebSocket server streams processed radar data in JSON format.

## Dependencies

- [gin-gonic/gin](https://github.com/gin-gonic/gin): Web framework
- [google/gopacket](https://github.com/google/gopacket): PCAP file processing
- [gorilla/websocket](https://github.com/gorilla/websocket): WebSocket implementation

## Project Structure

```
.
├── data/                   # PCAP data files
├── global/                 # Global variables and configurations
├── parser/                 # ASTERIX message parser
│   └── cat240/            # Category 240 specific parsing
├── sender/                # WebSocket server
├── utils/                 # Utility functions
├── docker-compose.yml     # Docker compose configuration
├── dockerfile             # Docker build configuration
├── go.mod                 # Go module file
└── main.go               # Application entry point
```

