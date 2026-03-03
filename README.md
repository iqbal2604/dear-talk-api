# DearTalk API

DearTalk is a real-time chat application backend built with Go, Gin, WebSocket, PostgreSQL, and Redis. It provides comprehensive APIs for user authentication, private and group chat rooms, and real-time messaging.

## Tech Stack
- **Go**: `1.25.5`
- **Framework**: [Gin](https://gin-gonic.com/)
- **Database**: PostgreSQL with [GORM](https://gorm.io/)
- **Cache / Rate Limiting**: Redis
- **Real-time**: Gorilla WebSocket
- **Dependency Injection**: Google Wire
- **Security**: JWT & bcrypt

## Project Structure
- `cmd/`: Application entry points.
- `internal/`: Private application and library code.
  - `domain/`: Business entities and models.
  - `router/`: Route definitions and handlers grouping.
  - `handler/`: HTTP request handlers.
  - `usecase/`: Business logic layer.
  - `repository/`: Data access layer.
  - `middleware/`: HTTP middlewares (Auth, CORS, Rate Limiting).
  - `websocket/`: Real-time WebSocket handlers and hub.
- `pkg/`: Publicly usable libraries (e.g. response formatters).
- `apispec.md`: Comprehensive API specifications.

## Prerequisites
- Go 1.25+
- PostgreSQL
- Redis

## Setup & Run

1. **Clone the repository:**
   ```bash
   git clone https://github.com/iqbal2604/dear-talk-api.git
   cd dear-talk-api
   ```

2. **Environment Variables:**
   Copy the example environment file and update with your local credentials:
   ```bash
   cp .env.example .env
   ```

3. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the Application:**
   ```bash
   go run cmd/main.go
   ```
   The server will start on `http://localhost:8080` (or your configured `APP_PORT`).

## API Documentation
Please refer to [apispec.md](./apispec.md) for the detailed API endpoint routes, methods, and responsibilities.

## Features
- **Authentication**: JWT-based user signup, login, and logout.
- **User Management**: Profile viewing, updating, and user search.
- **Rooms**: Support for 1-to-1 conversations and group chats, allowing adding/removing members.
- **Real-time Messaging**: WebSocket endpoints for live chat updates.
- **Rate Limiting**: Redis-backed rate limiter for application endpoints.

## License
MIT
