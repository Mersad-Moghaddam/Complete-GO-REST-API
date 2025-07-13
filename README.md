# GO Gin REST API

A complete RESTful API built with Go, Gin, and SQLite.  
Includes JWT authentication, event management, attendee management, and auto-generated Swagger documentation.

## Features

- User registration & login (JWT authentication)
- CRUD for events
- Manage event attendees
- SQLite database with migrations
- Auto-generated Swagger UI (`/swagger`)

## Getting Started

### Prerequisites

- Go 1.18+
- SQLite3

### Setup

1. **Clone the repository:**
   ```sh
   git clone https://github.com/your-username/your-repo.git
   cd your-repo
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Run database migrations:**
   ```sh
   go run cmd/migrate/main.go up
   ```

4. **Generate Swagger docs:**
   ```sh
   swag init
   ```

5. **Start the API server:**
   ```sh
   go run cmd/api/main.go
   ```

6. **Access Swagger UI:**
   - [http://localhost:8080/swagger](http://localhost:8080/swagger)
   - [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## API Endpoints

See the Swagger UI for full documentation.

## Environment Variables

- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: Secret for JWT signing (default: "secret")

## Project Structure

```
cmd/api/         # Main API server
cmd/migrate/     # Database migration tool
internal/        # Application logic (database, env)
docs/            # Swagger docs (auto-generated)
```

## License

MIT

