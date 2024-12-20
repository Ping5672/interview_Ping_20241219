# OXO Game Service Backend

## Overview

The OXO Game Service Backend is a comprehensive game management system built to provide robust player interaction, game room management, challenging gameplay, and secure payment processing.

## Features

### 1. Player Management
- Create and manage player profiles
- Implement player leveling system
- Perform CRUD operations on player data

### 2. Game Room Management
- Create and manage game rooms
- Advanced reservation system
- Real-time room status tracking

### 3. Endless Challenge System
- 30-second challenge duration
- Exciting 1% win probability
- Fixed 20.01 entry fee
- Dynamic prize pool accumulation

### 4. Game Log Collector
- Comprehensive action logging
- Detailed player activity tracking
- Flexible log retrieval with filtering

### 5. Payment Processing
- Support for multiple payment methods
- Secure transaction tracking
- Comprehensive payment processing

## Technology Stack

- **Language**: Go (1.21+)
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Caching**: Redis
- **Containerization**: Docker

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL
- Redis

## Quick Start

### Installation

1. Clone the repository:
```bash
git clone <repository_url>
cd interview_Ping_20241219
```

2. Build and start services:
```bash
docker-compose up --build -d
```

3. Verify services are running:
```bash
docker-compose ps
```

## API Endpoints

### Player Management
- **Create Player**: `POST /players`
- **List Players**: `GET /players`
- **Get Player**: `GET /players/{id}`

### Room Management
- **Create Room**: `POST /rooms`
- **List Rooms**: `GET /rooms`
- **Make Reservation**: `POST /reservations`

### Challenge System
- **Join Challenge**: `POST /challenges`
- **Get Challenge Results**: `GET /challenges/results`

### Game Logging
- **Create Log**: `POST /logs`
- **Retrieve Logs**: `GET /logs` (with optional filtering)

### Payment Processing
- **Process Payment**: `POST /payments`
- **Check Payment Status**: `GET /payments/{id}`

## Payment Method Details

### Credit Card
- Processing Time: 800ms
- Failure Rate: 10%
- Transaction ID Format: `CC_CARD_*`

### Bank Transfer
- Processing Time: 1s
- Failure Rate: 5%
- Transaction ID Format: `BT_BANK_*`

### Third Party Payment
- Processing Time: 600ms
- Failure Rate: 8%
- Transaction ID Format: `TP_3RDPARTY_*`

### Blockchain Payment
- Processing Time: 2s
- Failure Rate: 15%
- Transaction ID Format: `BC_*`

## Error Handling

### Common Error Codes
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `429`: Too Many Requests
- `500`: Internal Server Error

### Error Response Format
```json
{
    "error": "Error message",
    "details": "Detailed error information",
    "code": "Error code (if applicable)"
}
```

## Testing

Run all tests:
```bash
go test ./tests/... -v
```

Run specific tests:
```bash
go test ./tests/payment_test.go -v
go test ./tests/room_test.go -v
```

## Monitoring and Maintenance

### Check Service Status
```bash
docker-compose ps
```

### View Logs
```bash
docker-compose logs -f app
docker-compose logs -f db
docker-compose logs -f redis
```

### Database Maintenance
```bash
docker exec -it interview_ping_20241219-db-1 psql -U postgres -d game_db
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit changes
4. Push to the branch
5. Create a Pull Request
