# 🚗 SpotSync API

SpotSync API is a Go backend for smart parking and EV charging reservation management. It provides authentication, role-based access control, parking zone management, live spot availability, and concurrency-safe reservation handling so the final available slot cannot be double-booked.

## Live Project

- Repository: [PH-Level-2-Assignment-6_Spot-Sync](https://github.com/SamiunAuntor/PH-Level-2-Assignment-6_Spot-Sync)
- Live API: [https://spot-sync-api.onrender.com/](https://spot-sync-api.onrender.com/)
- API Testing Guide: [postman/API_TESTING.md](https://github.com/SamiunAuntor/PH-Level-2-Assignment-6_Spot-Sync/blob/main/postman/API_TESTING.md)

## Overview

This project was built as a backend-only parking reservation system where:

- users can register and log in with JWT authentication
- drivers can browse parking zones, reserve spots, and manage their own reservations
- admins can manage parking zones and inspect all reservations
- live availability is derived from active reservations
- reservation creation is protected with transactions and row-level locking

## Core Features

- JWT-based authentication for `driver` and `admin`
- Clean layered architecture with manual dependency injection
- Public parking zone listing and single-zone retrieval
- Admin parking zone create, update, and delete
- Authenticated reservation creation
- Owner-only reservation cancellation
- Admin access to all reservations
- Centralized validation and error handling
- PostgreSQL persistence with GORM
- Concurrency-safe spot booking using `FOR UPDATE`

## Tech Stack

| Technology | Purpose |
| --- | --- |
| Go 1.22+ | Backend language |
| Echo v4 | HTTP framework |
| GORM | Database ORM |
| PostgreSQL | Primary database |
| JWT v5 | Authentication |
| bcrypt | Password hashing |
| validator/v10 | Request validation |
| Render | Backend deployment |
| NeonDB | Hosted PostgreSQL |

## Architecture

The codebase follows a clean, layered backend structure:

```text
DTO -> Handler -> Service -> Repository -> GORM -> PostgreSQL
```

Layer responsibilities:

- `dto/` defines request payloads and transport-facing structures
- `handler/` parses requests and returns HTTP responses
- `service/` contains business rules and orchestration
- `repository/` handles database access, transactions, and locking
- `models/` defines persistent entities and GORM mappings
- `middleware/` handles auth, roles, and centralized error flow

Dependency wiring is done manually in `main.go`.

## Roles and Permissions

### Driver

- Register and log in
- View all parking zones and availability
- Create reservations
- View personal reservations
- Cancel personal reservations

### Admin

- All driver capabilities
- Create parking zones
- Update parking zones
- Delete parking zones
- View all reservations

## Concurrency Safety

The reservation flow protects against overbooking when multiple users try to claim the final remaining spot at the same time.

How it works:

- a database transaction is started
- the selected parking zone row is locked with `FOR UPDATE`
- active reservations for that zone are counted inside the same transaction
- capacity is validated before insertion
- the reservation is created atomically only if space still exists

This ensures a zone never exceeds its configured `total_capacity`.

## Environment Variables

Use `.env.example` as the reference.

| Variable | Required | Description |
| --- | --- | --- |
| `PORT` | No | Server port, defaults to `8080` locally |
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | JWT signing secret |
| `JWT_EXPIRES_IN` | No | Token lifetime, default `24h` |
| `CORS_ALLOWED_ORIGINS` | No | Allowed origins, default `*` |

## Database Schema

### `users`

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `int` | Primary key, auto-increment |
| `name` | `text` | Required |
| `email` | `text` | Required, unique |
| `password` | `text` | Hashed with bcrypt |
| `role` | `varchar(20)` | `driver` or `admin` |
| `created_at` | `timestamptz` | Auto managed |
| `updated_at` | `timestamptz` | Auto managed |

### `parking_zones`

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `int` | Primary key, auto-increment |
| `name` | `text` | Required |
| `type` | `varchar(30)` | `general`, `ev_charging`, or `covered` |
| `total_capacity` | `integer` | Must be greater than `0` |
| `price_per_hour` | `numeric(10,2)` | Must be greater than `0` |
| `created_at` | `timestamptz` | Auto managed |
| `updated_at` | `timestamptz` | Auto managed |

### `reservations`

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `int` | Primary key, auto-increment |
| `user_id` | `int` | Foreign key to `users.id` |
| `zone_id` | `int` | Foreign key to `parking_zones.id` |
| `license_plate` | `varchar(15)` | Required |
| `status` | `varchar(20)` | `active`, `completed`, or `cancelled` |
| `created_at` | `timestamptz` | Auto managed |
| `updated_at` | `timestamptz` | Auto managed |

## Common HTTP Status Codes

| Status Code | Meaning | Typical Usage In This Project |
| --- | --- | --- |
| `200 OK` | Request succeeded | Fetching resources, successful cancel/update |
| `201 Created` | Resource created successfully | Registration, zone creation, reservation creation |
| `400 Bad Request` | Invalid input or validation failure | Invalid body, malformed ID, bad request data |
| `401 Unauthorized` | Missing or invalid authentication | Missing token, invalid token, failed login |
| `403 Forbidden` | Authenticated but not allowed | Driver trying admin-only endpoint |
| `404 Not Found` | Requested resource does not exist | Zone or reservation not found |
| `409 Conflict` | State conflict with business rules | Zone full, duplicate plate, invalid reservation state |
| `500 Internal Server Error` | Unexpected server-side error | Unhandled internal failure |

## Getting Started

### 1. Clone the repository

```powershell
git clone https://github.com/SamiunAuntor/PH-Level-2-Assignment-6_Spot-Sync.git
cd PH-Level-2-Assignment-6_Spot-Sync
```

### 2. Configure environment variables

Create a `.env` file based on `.env.example` and set at least:

- `DATABASE_URL`
- `JWT_SECRET`

### 3. Install dependencies

```powershell
go mod tidy
```

### 4. Run the server

```powershell
go run .
```

Local base URL:

```text
http://localhost:8080
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

### Parking Zones

- `POST /api/v1/zones`
- `GET /api/v1/zones`
- `GET /api/v1/zones/:id`
- `PATCH /api/v1/zones/:id`
- `DELETE /api/v1/zones/:id`

### Reservations

- `POST /api/v1/reservations`
- `GET /api/v1/reservations/my-reservations`
- `DELETE /api/v1/reservations/:id`
- `GET /api/v1/reservations`

## Example Requests

### Register User

```json
{
  "name": "John Doe",
  "email": "john.doe@spotsync.com",
  "password": "securePassword123",
  "role": "driver"
}
```

### Login

```json
{
  "email": "john.doe@spotsync.com",
  "password": "securePassword123"
}
```

### Create Parking Zone

```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.5
}
```

### Create Reservation

```json
{
  "zone_id": 1,
  "license_plate": "ABC-1234"
}
```

For full endpoint-by-endpoint request and response examples, see the API testing guide:

- [API_TESTING.md](https://github.com/SamiunAuntor/PH-Level-2-Assignment-6_Spot-Sync/blob/main/postman/API_TESTING.md)

## Project Structure

```text
apperror/
config/
database/
dto/
handler/
middleware/
models/
postman/
repository/
response/
routes/
service/
validator/
.env.example
main.go
README.md
```

## Notes

- `CORS_ALLOWED_ORIGINS=*` is acceptable for testing across different machines
- for production frontend integration, restrict CORS to trusted origins
- the extra zone update and delete endpoints use the same `/api/v1/zones/:id` URL pattern as the rest of the API
