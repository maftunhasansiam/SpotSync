# SpotSync API

> Backend API for managing parking zones and EV charging reservations at high-traffic locations like airports and malls.

Built with **Go** · **Echo v4** · **GORM** · **PostgreSQL**

🔗 Live URL: `https://spotsync-12.onrender.com`

---

## Features

- JWT authentication with **Driver** and **Admin** roles
- Admins can create, update, and delete parking zones
- Drivers can browse zones and see live available spots
- Drivers can reserve spots — admins can view all reservations
- **Concurrency-safe booking** using PostgreSQL row-level locking — no over-capacity, ever

---

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go 1.22 |
| HTTP Framework | Echo v4 |
| ORM | GORM |
| Database | PostgreSQL (NeonDB recommended) |
| Auth | JWT (`golang-jwt/jwt v5`) + bcrypt cost 12 |
| Validation | `go-playground/validator v10` |
| Hot Reload | Air |

---

## Architecture

Follows strict **Clean Architecture** — each layer has one responsibility.

```
HTTP Request
     ↓
  Handler        → binds request, validates input, calls service, returns JSON
     ↓
  Service        → business logic (hashing, JWT, capacity checks)
     ↓
  Repository     → all database calls (GORM, transactions, row locks)
     ↓
  PostgreSQL
```

Dependency injection is wired manually in `main.go`:
`repository → service → handler`

Handlers never touch the database directly.

---

## Project Structure

```
.
├── main.go                  # Entry point, DI, route registration
├── config/                  # Database connection + AutoMigrate
├── models/                  # GORM structs (User, ParkingZone, Reservation)
├── dto/                     # Request and response payloads
├── middleware/              # JWT verification + AdminOnly guard
├── repository/              # All DB operations
├── service/                 # Business logic
├── handler/                 # HTTP layer
└── handler/validation.go    # Shared validator error helper
```

---

## Getting Started

**Prerequisites:** Go 1.22+, a PostgreSQL database (NeonDB works great)

```bash
# Clone the repo
git clone https://github.com/maftunhasansiam/SpotSync
cd SpotSync

# Install dependencies
go mod tidy

# Copy env template and fill in your values
cp .env.example .env

# Run with hot reload (recommended)
go install github.com/air-verse/air@latest
air

# Or run directly
go run main.go
```

Server starts on port `8080` by default. Set `PORT` in `.env` to change.

---

## Environment Variables

```env
DATABASE_URL=postgres://user:***@host/dbname?sslmode=require
JWT_SECRET=replace-with-a-long-random-string
PORT=8080
```

---

## API Endpoints

### Auth — Public

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Log in, returns JWT |

### Parking Zones

| Method | Endpoint | Access | Description |
|---|---|---|---|
| GET | `/api/v1/zones` | Public | List all zones with available spots |
| GET | `/api/v1/zones/:id` | Public | Get a single zone |
| POST | `/api/v1/zones` | Admin | Create a new zone |
| PUT | `/api/v1/zones/:id` | Admin | Update a zone |
| DELETE | `/api/v1/zones/:id` | Admin | Delete a zone |

### Reservations

| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | `/api/v1/reservations` | Authenticated | Reserve a spot |
| GET | `/api/v1/reservations/my-reservations` | Authenticated | List my reservations |
| DELETE | `/api/v1/reservations/:id` | Authenticated | Cancel my reservation |
| GET | `/api/v1/reservations` | Admin | List all reservations |

> Send the JWT in the `Authorization: Bearer <token>` header for all protected routes.

---

## Notes

- **Concurrency safety** — The reservation endpoint uses a PostgreSQL transaction with `FOR UPDATE` row-level locking on the parking zone, preventing over-booking under concurrent load.
- **Password security** — Passwords are stored as bcrypt hashes (cost 12) and are never returned in any API response.
- **Available spots** — The `available_spots` field on zones is calculated at read time, not stored, so it always reflects the real current state.
