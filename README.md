# SpotSync API

Backend API for managing parking zones and EV charging reservations at high-traffic locations like airports and malls. Built with Go + Echo + GORM + PostgreSQL.

Live URL: https://spotsync-api.onrender.com  *(will be update after deployment)*

---

## What it does

- JWT authentication with driver/admin roles
- Admin can create, update, and delete parking zones
- Drivers can browse zones and see live available spots
- Drivers can reserve spots; admins can view all reservations
- Concurrency-safe booking using Postgres row-level locking (no over-capacity)

---

## Tech stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.22 |
| HTTP framework | Echo v4 |
| ORM | GORM |
| Database | PostgreSQL (NeonDB recommended) |
| Auth | JWT (golang-jwt/jwt v5) + bcrypt cost 12 |
| Validation | go-playground/validator v10 |
| Hot reload | Air |

---

## Architecture

The project follows strict Clean Architecture. Each layer has one job:
HTTP request | v Handler -- binds request, validates, calls service, returns JSON | v Service -- business logic (hashing, JWT, capacity checks) | v Repository -- all database calls (GORM, transactions, row locks) | v PostgreSQL

Code
Dependency injection is wired manually in `main.go`: repository -> service -> handler. Handlers never touch the database directly.

---

## Project layout
. ├── main.go # Entry point, DI, route registration ├── config/ # Database connection + AutoMigrate ├── models/ # GORM structs (User, ParkingZone, Reservation) ├── dto/ # Request and response payloads ├── middleware/ # JWT verification + AdminOnly guard ├── repository/ # All DB operations ├── service/ # Business logic ├── handler/ # HTTP layer └── handler/validation.go # Shared validator error helper

Code
---

## Running locally

Prerequisites: Go 1.22+, a PostgreSQL database (NeonDB works).

```bash
git clone <https://github.com/maftunhasansiam/SpotSync>
cd SpotSync

# Install dependencies
go mod tidy

# Copy env template and fill in your values
cp .env.example .env
# Edit .env with your DATABASE_URL and JWT_SECRET

# Run with hot reload (recommended for dev)
go install github.com/air-verse/air@latest
air

# Or run directly
go run main.go
The server starts on port 8080 by default (set PORT in .env to change).

Required env variables
DATABASE_URL=postgres://user:***@host/dbname?sslmode=require
JWT_SECRET=replace-with-a-long-random-string
PORT=8080
API endpoints
Auth (Public)
POST /api/v1/auth/register — register a new user
POST /api/v1/auth/login — log in, returns JWT
Parking zones
GET /api/v1/zones — list all zones with available spots (Public)
GET /api/v1/zones/:id — get one zone (Public)
POST /api/v1/zones — create a zone (Admin)
PUT /api/v1/zones/:id — update a zone (Admin)
DELETE /api/v1/zones/:id — delete a zone (Admin)
Reservations
POST /api/v1/reservations — reserve a spot (Authenticated)
GET /api/v1/reservations/my-reservations — list my reservations (Authenticated)
DELETE /api/v1/reservations/:id — cancel my reservation (Authenticated)
GET /api/v1/reservations — list all reservations (Admin)
Send the JWT in the Authorization: Bearer *** header for protected routes.

Notes
The reservation create endpoint uses a transaction + FOR UPDATE row lock on the parking zone to prevent over-booking under concurrent load.
Passwords are stored as bcrypt hashes (cost 12). They are never returned in API responses.
The available_spots field on zones is calculated at read time, not stored.