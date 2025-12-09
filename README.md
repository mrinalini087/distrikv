Distributed key-value storage system built in Go (Golang).
Demonstrates core distributed computing concepts including sharding, load balancing, and consistent hashing.
The system is composed of multiple independent storage nodes and a central load balancer (proxy) that routes requests to the correct node based on the data key.

## Key Features

- **Distributed Architecture:** Data is partitioned across multiple storage nodes (Sharding)
- **Load Balancing:** A central proxy distributes traffic efficiently
- **Consistent Hashing:** Uses a hashing algorithm to deterministically map keys to specific servers
- **Concurrency Safe:** Implements `sync.RWMutex` to handle concurrent read/write operations safely
- **REST API:** Simple HTTP interface for setting and retrieving values

## Tech Stack

- **Language:** Go (Golang) 1.23+
- **Communication:** HTTP / JSON
- **Architecture:** Client-Server / Microservices

## Quick Start

**4 separate terminal windows** to run the distributed system
Split terminal into 4
Run each command in a separate terminal:

```bash
# Terminal 1 - Node A
go run node/main.go -port=8081
```

```bash
# Terminal 2 - Node B
go run node/main.go -port=8082
```

```bash
# Terminal 3 - Node C
go run node/main.go -port=8083
```

### Start the Proxy Server

```bash
# Terminal 4 - Proxy
go run proxy/main.go
```

The proxy will automatically distribute requests across all three storage nodes using consistent hashing.
Once the system is running, interact with it using HTTP requests:

### Set a Key-Value Pair

Using curl commands on a separate terminal, we can see if it works or not

```bash
curl -X POST http://localhost:8080/set \
  -H "Content-Type: application/json" \
  -d '{"key": "username", "value": "alice"}'
```

### Get a Value

```bash
curl http://localhost:8080/get?key=username
```

## Architecture Overview

```
Client Request
      ↓
  Proxy Server (Port 8080)
      ↓
  Consistent Hashing
      ↓
  ┌─────────┬─────────┬─────────┐
  Node A    Node B    Node C
  :8081     :8082     :8083
```

The proxy uses consistent hashing to determine which storage node should handle each key, ensuring balanced distribution and predictable routing.

## Requirements

- Go 1.23 or higher
- 4 available ports (8080-8083)

This project was for learning purposes, the aim was for:

- Building distributed systems in Go
- Implementing consistent hashing algorithms
- Handling concurrent operations with mutexes
- Creating REST APIs with Go's standard library
- Microservices communication patterns

## License

MIT License - feel free to use this project for learning and experimentation.
