# Go Real-Time Chat

A lightweight, room-based real-time chat application built with **Go and WebSockets**. The application enables multiple clients to communicate through persistent WebSocket connections, with room-based message broadcasting and concurrent client management.

## Features

* **Real-time Messaging:** Exchange messages instantly using WebSocket connections.
* **Room-Based Communication:** Users can join rooms and communicate with other users in the same room.
* **Concurrent Client Management:** Handle multiple connected clients using Go's concurrency primitives.
* **Username Validation:** Prevent duplicate usernames among connected users.
* **System Messages:** Broadcast notifications when users leave a room.
* **Structured JSON Messages:** Send messages containing the message type, sender, content, and timestamp.
* **Connection Management:** Register and unregister clients as connections are established and terminated.
* **Simple Web Interface:** Connect to a room, send messages, and view incoming messages.

## Tech Stack

| Technology        | Purpose                              |
| ----------------- | ------------------------------------ |
| Go                | Backend server                       |
| Gorilla WebSocket | WebSocket communication              |
| HTML              | Frontend structure                   |
| JavaScript        | WebSocket client and UI interactions |

## Architecture

The application follows a Hub-and-Client architecture.

* **Hub:** Maintains connected clients and coordinates message broadcasting.
* **Client:** Represents an individual WebSocket connection.
* **Broadcast Channel:** Receives messages and distributes them to clients in the corresponding room.
* **Read/Write Pumps:** Handle WebSocket communication and outgoing message delivery.

### Message Flow

1. A client connects to the server through the WebSocket endpoint.
2. The server validates the username and registers the client with the Hub.
3. The client sends a message.
4. The server broadcasts the message to connected clients in the same room.
5. Messages are delivered as structured JSON objects.

## Project Structure

```text
go-chat/
├── .gitignore
├── go.mod
├── go.sum
├── hub.go
├── index.html
├── main.go
└── README.md

## Getting Started

### Prerequisites

* Go 1.25 or compatible version
* Git

### Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/go-chat.git
cd go-chat
```

### Install Dependencies

```bash
go mod tidy
```

### Run the Server

```bash
go run .
```

The server runs at:

```text
http://localhost:8080
```

Open the URL in your browser to access the chat interface.

## WebSocket Endpoint

```text
ws://localhost:8080/ws
```

### Connection Parameters

| Parameter  | Description                   |
| ---------- | ----------------------------- |
| `username` | Username chosen by the client |
| `room`     | Chat room to join             |

Example:

```text
ws://localhost:8080/ws?username=Shaurya&room=general
```

## Message Format

Messages are exchanged in JSON format.

```json
{
  "message": "Hello everyone!",
  "type": "chat",
  "sender": "Shaurya",
  "timestamp": "10:30 PM"
}
```

System events use the `system` message type.

