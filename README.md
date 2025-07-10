# Chirpy-Go
Chirpy-Go is a simple social media API server written in Go. It provides endpoints for user management, authentication, posting short messages ("chirps"), and basic admin functionality. The server uses PostgreSQL for data storage and supports JWT-based authentication.

## Features
- User registration and update
- JWT-based login, refresh, and revoke
- Posting, retrieving, and deleting chirps
- Webhook support for Polka
- Admin endpoints for metrics and reset
- File server for static assets

## Endpoints
### Public Endpoints
- `GET /api/healthz` — Health check
- `POST /api/login` — User login (JWT)
- `POST /api/refresh` — Refresh JWT token
- `POST /api/revoke` — Revoke JWT token
- `POST /api/users` — Create a new user
- `PUT /api/users` — Update user info
- `GET /api/chirps` — List all chirps
- `GET /api/chirps/{chirpID}` — Get a specific chirp
- `POST /api/chirps` — Create a new chirp
- `DELETE /api/chirps/{chirpID}` — Delete a chirp
- `POST /api/polka/webhooks` — Handle Polka webhooks
  
### Admin Endpoints
- `POST /admin/reset` — Reset the application state
- `GET /admin/metrics` — Get server metrics

### Static Files
- `/app/` — Serves static files from the project root

## Configuration
The server uses environment variables for configuration:

- `DB_URL` — PostgreSQL connection string (required)
- `PLATFORM` — Platform identifier (required)
- `JWT_SECRET` — Secret for signing JWT tokens (required)
- `POLKA_KEY` — Key for Polka webhook validation (required)
  
You can use a .env file for local development. The server loads environment variables using [joho/godotenv](https://github.com/joho/godotenv).

## Running the Server
1. Install dependencies:
   ```go
   go mod tidy
   ```
2. Set up your .env file with the required variables.
3. Start the server:
   ```go
   go run main.go
   ```
The server will listen on port `8080` by default.

## Project Structure
- main.go — Entry point and HTTP server setup
- internal/database/ — Database access layer
- handler_*.go — HTTP handlers for various endpoints
- index.html — Example static file
