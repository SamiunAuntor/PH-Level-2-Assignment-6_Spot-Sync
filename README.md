# SpotSync API

SpotSync is a Go backend for smart parking and EV charging reservation management. It supports user authentication, parking zone management, live availability tracking, and reservation handling with concurrency-safe booking logic to prevent overbooking limited spots.

## Live URL

```text
TBD
```

## Features

- User registration and login with JWT authentication
- Role-based access control for `driver` and `admin`
- Public parking zone listing with dynamic `available_spots`
- Admin parking zone creation, update, and deletion
- Reservation creation for authenticated users
- Owner-only reservation cancellation
- Admin access to all reservation records
- Concurrency-safe reservation creation using GORM transactions and row locking

## Tech Stack

| Technology | Usage |
| --- | --- |
| Go 1.22+ | Backend language |
| Echo v4 | HTTP framework |
| GORM | ORM |
| PostgreSQL | Database |
| validator/v10 | Request validation |
| JWT v5 | Authentication |
| bcrypt | Password hashing |

## Architecture

The project follows the clean architecture rules required in the assignment:

```text
DTO -> Handler -> Service -> Repository -> GORM -> PostgreSQL
```

Responsibilities:

- `dto/`: request and response payloads
- `handler/`: HTTP request handling and response writing
- `service/`: business rules and orchestration
- `repository/`: database queries, transactions, and row locks
- `models/`: GORM entity definitions

Manual dependency injection is done in `main.go`.

## Roles and Permissions

### Driver

- Register and log in
- View all parking zones and availability
- Reserve a parking or EV charging spot
- View their own reservations
- Cancel their own reservations

### Admin

- All driver permissions
- Create parking zones
- Update parking zones
- Delete parking zones
- Set pricing through zone updates
- View all reservations

## Environment Variables

Use `.env.example` as the reference for local configuration.

| Variable | Required | Description |
| --- | --- | --- |
| `PORT` | No | Server port, default `8080` |
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | JWT signing secret |
| `JWT_EXPIRES_IN` | No | Token lifetime, default `24h` |
| `CORS_ALLOWED_ORIGINS` | No | Allowed origins, default `*` |

## Local Setup

### 1. Clone the repository

```powershell
git clone <repository-url>
cd SpotSync
```

### 2. Configure environment variables

Create a `.env` file and provide values for:

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

Default local URL:

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

### Register

```http
POST /api/v1/auth/register
Content-Type: application/json
```

```json
{
  "name": "John Doe",
  "email": "john.doe@spotsync.com",
  "password": "securePassword123",
  "role": "driver"
}
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "john.doe@spotsync.com",
  "password": "securePassword123"
}
```

### Create Parking Zone

```http
POST /api/v1/zones
Authorization: Bearer <admin_token>
Content-Type: application/json
```

```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.5
}
```

### Create Reservation

```http
POST /api/v1/reservations
Authorization: Bearer <driver_or_admin_token>
Content-Type: application/json
```

```json
{
  "zone_id": 1,
  "license_plate": "ABC-1234"
}
```

## Concurrency Handling

Reservation creation is protected from overbooking by:

- starting a database transaction
- locking the selected parking zone row with `FOR UPDATE`
- counting active reservations inside the same transaction
- rejecting booking when the zone is full
- inserting the reservation inside the same transaction

This ensures two simultaneous requests cannot both reserve the final available spot.

## CORS

`CORS_ALLOWED_ORIGINS=*` is acceptable for development and testing. For deployment, replace it with the actual frontend origin or a comma-separated list of trusted origins.

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
main.go
```

## Project Links

### Repository

```text
TBD
```

### Live Deployment

```text
TBD
```

### Interview Video

```text
TBD
```
