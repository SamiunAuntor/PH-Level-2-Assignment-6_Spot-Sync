# SpotSync API Testing Guide

This file is a quick manual testing guide and lightweight API reference for all implemented endpoints.

## Base URL

```text
http://localhost:8080
```

## Auth Tokens

Use the login endpoint first and copy the returned JWT token.

Example auth header:

```text
Authorization: Bearer <your_jwt_token>
```

## Roles Used In Examples

- `driver`
- `admin`

## 1. Register User

**Endpoint**

```text
POST /api/v1/auth/register
```

**Access**

Public

**Expected Input**

```json
{
  "name": "John Doe",
  "email": "john.doe@spotsync.com",
  "password": "securePassword123",
  "role": "driver"
}
```

**Expected Success**

- Status: `201`
- Message: `User registered successfully`

Example response:

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@spotsync.com",
    "role": "driver",
    "created_at": "2026-06-26T10:00:00Z",
    "updated_at": "2026-06-26T10:00:00Z"
  }
}
```

**Negative Checks**

- duplicate email -> `400`
- invalid email -> `400`
- password shorter than 8 characters -> `400`

## 2. Login

**Endpoint**

```text
POST /api/v1/auth/login
```

**Access**

Public

**Expected Input**

```json
{
  "email": "john.doe@spotsync.com",
  "password": "securePassword123"
}
```

**Expected Success**

- Status: `200`
- Message: `Login successful`

Example response:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "your_jwt_here",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@spotsync.com",
      "role": "driver"
    }
  }
}
```

**Negative Checks**

- wrong password -> `401`
- unknown email -> `401`

## 3. Create Parking Zone

**Endpoint**

```text
POST /api/v1/zones
```

**Access**

Admin only

**Expected Input**

```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.5
}
```

**Expected Success**

- Status: `201`
- Message: `Parking zone created successfully`

Example response:

```json
{
  "success": true,
  "message": "Parking zone created successfully",
  "data": {
    "id": 1,
    "name": "Terminal 1 EV Charging",
    "type": "ev_charging",
    "total_capacity": 20,
    "price_per_hour": 5.5,
    "created_at": "2026-06-26T10:10:00Z",
    "updated_at": "2026-06-26T10:10:00Z"
  }
}
```

**Negative Checks**

- driver token -> `403`
- missing token -> `401`
- invalid type -> `400`
- zero or negative capacity -> `400`
- zero or negative price -> `400`

## 4. Get All Parking Zones

**Endpoint**

```text
GET /api/v1/zones
```

**Access**

Public

**Expected Success**

- Status: `200`
- Message: `Parking zones retrieved successfully`

Example response:

```json
{
  "success": true,
  "message": "Parking zones retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Terminal 1 EV Charging",
      "type": "ev_charging",
      "total_capacity": 20,
      "available_spots": 19,
      "price_per_hour": 5.5,
      "created_at": "2026-06-26T10:10:00Z"
    }
  ]
}
```

## 5. Get Single Parking Zone

**Endpoint**

```text
GET /api/v1/zones/:id
```

**Access**

Public

**Example**

```text
GET /api/v1/zones/1
```

**Expected Success**

- Status: `200`
- Message: `Parking zone retrieved successfully`

Example response:

```json
{
  "success": true,
  "message": "Parking zone retrieved successfully",
  "data": {
    "id": 1,
    "name": "Terminal 1 EV Charging",
    "type": "ev_charging",
    "total_capacity": 20,
    "available_spots": 19,
    "price_per_hour": 5.5,
    "created_at": "2026-06-26T10:10:00Z"
  }
}
```

**Negative Checks**

- invalid numeric id -> `400`
- missing zone -> `404`

## 6. Update Parking Zone

**Endpoint**

```text
PATCH /api/v1/zones/:id
```

**Access**

Admin only

**Example Input**

```json
{
  "name": "Terminal 1 Covered Parking",
  "price_per_hour": 6.25
}
```

**Expected Success**

- Status: `200`
- Message: `Parking zone updated successfully`

Example response:

```json
{
  "success": true,
  "message": "Parking zone updated successfully",
  "data": {
    "id": 1,
    "name": "Terminal 1 Covered Parking",
    "type": "ev_charging",
    "total_capacity": 20,
    "price_per_hour": 6.25,
    "created_at": "2026-06-26T10:10:00Z",
    "updated_at": "2026-06-26T10:20:00Z"
  }
}
```

**Negative Checks**

- driver token -> `403`
- missing token -> `401`
- empty body -> `400`
- invalid zone id -> `400`
- missing zone -> `404`

## 7. Delete Parking Zone

**Endpoint**

```text
DELETE /api/v1/zones/:id
```

**Access**

Admin only

**Expected Success**

- Status: `200`
- Message: `Parking zone deleted successfully`

Example response:

```json
{
  "success": true,
  "message": "Parking zone deleted successfully"
}
```

**Negative Checks**

- driver token -> `403`
- missing token -> `401`
- missing zone -> `404`
- zone with reservation history -> `409`

## 8. Create Reservation

**Endpoint**

```text
POST /api/v1/reservations
```

**Access**

Authenticated `driver` or `admin`

**Expected Input**

```json
{
  "zone_id": 1,
  "license_plate": "ABC-1234"
}
```

**Expected Success**

- Status: `201`
- Message: `Reservation confirmed successfully`

Example response:

```json
{
  "success": true,
  "message": "Reservation confirmed successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "zone_id": 1,
    "license_plate": "ABC-1234",
    "status": "active",
    "created_at": "2026-06-26T10:30:00Z",
    "updated_at": "2026-06-26T10:30:00Z"
  }
}
```

**Negative Checks**

- missing token -> `401`
- invalid zone id -> `400`
- missing zone -> `404`
- zone full -> `409`
- duplicate active license plate -> `409`

## 9. Get My Reservations

**Endpoint**

```text
GET /api/v1/reservations/my-reservations
```

**Access**

Authenticated `driver` or `admin`

**Expected Success**

- Status: `200`
- Message: `My reservations retrieved successfully`

Example response:

```json
{
  "success": true,
  "message": "My reservations retrieved successfully",
  "data": [
    {
      "id": 1,
      "license_plate": "ABC-1234",
      "status": "active",
      "zone": {
        "id": 1,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "created_at": "2026-06-26T10:30:00Z"
    }
  ]
}
```

## 10. Cancel Reservation

**Endpoint**

```text
DELETE /api/v1/reservations/:id
```

**Access**

Authenticated owner only

**Expected Success**

- Status: `200`
- Message: `Reservation cancelled successfully`

Example response:

```json
{
  "success": true,
  "message": "Reservation cancelled successfully"
}
```

**Negative Checks**

- missing token -> `401`
- another user's reservation -> `403`
- missing reservation -> `404`
- already cancelled/completed -> `409`

## 11. Get All Reservations

**Endpoint**

```text
GET /api/v1/reservations
```

**Access**

Admin only

**Expected Success**

- Status: `200`
- Message: `All reservations retrieved successfully`

Example response:

```json
{
  "success": true,
  "message": "All reservations retrieved successfully",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "zone_id": 1,
      "license_plate": "ABC-1234",
      "status": "active",
      "user": {
        "id": 1,
        "name": "John Doe",
        "email": "john.doe@spotsync.com",
        "role": "driver"
      },
      "zone": {
        "id": 1,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "created_at": "2026-06-26T10:30:00Z",
      "updated_at": "2026-06-26T10:30:00Z"
    }
  ]
}
```

**Negative Checks**

- driver token -> `403`
- missing token -> `401`

## Suggested Test Order

1. Register a driver
2. Register an admin
3. Login as driver
4. Login as admin
5. Create a parking zone as admin
6. List and get zones publicly
7. Create a reservation as driver
8. Get my reservations as driver
9. Try to get all reservations as driver and confirm `403`
10. Get all reservations as admin
11. Cancel the reservation as driver
12. Try cancelling again and confirm `409`
