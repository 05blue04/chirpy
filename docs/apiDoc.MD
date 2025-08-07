# Chirpy API Documentation

A Twitter-like social media API for creating and managing short messages called "chirps".

## Base URL
```
http://localhost:8080
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

Some endpoints require API key authentication via the `Authorization` header:
```
Authorization: ApiKey <your_api_key>
```

## Users

### Create User
Create a new user account.

**Endpoint:** `POST /api/users`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "your_password"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### Login
Authenticate a user and receive access and refresh tokens.

**Endpoint:** `POST /api/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "your_password",
  "expires_in_seconds": 3600
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "random_refresh_token_string",
  "is_chirpy_red": false
}
```

### Update User
Update user information (requires authentication).

**Endpoint:** `PUT /api/users`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "password": "new_password"
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

## Authentication Tokens

### Refresh Token
Get a new access token using a refresh token.

**Endpoint:** `POST /api/refresh`

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Revoke Token
Revoke a refresh token.

**Endpoint:** `POST /api/revoke`

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Response:** `204 No Content`

## Chirps

### Create Chirp
Create a new chirp (requires authentication).

**Endpoint:** `POST /api/chirps`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "body": "This is my first chirp!"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "body": "This is my first chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Notes:**
- Maximum chirp length: 140 characters
- Profanity filter: The words "kerfuffle", "sharbert", and "fornax" will be replaced with "****"

### Get All Chirps
Retrieve all chirps.

**Endpoint:** `GET /api/chirps`

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z",
    "body": "This is my first chirp!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
]
```

### Get Chirp by ID
Retrieve a specific chirp by its ID.

**Endpoint:** `GET /api/chirps/{chirpID}`

**Path Parameters:**
- `chirpID`: UUID of the chirp

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "body": "This is my first chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Delete Chirp
Delete a chirp (requires authentication and ownership).

**Endpoint:** `DELETE /api/chirps/{chirpID}`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `chirpID`: UUID of the chirp to delete

**Response:** `204 No Content`

**Notes:**
- Only the author of the chirp can delete it

## Webhooks

### Polka Webhook
Handle webhook events from Polka payment system.

**Endpoint:** `POST /api/polka/webhooks`

**Headers:**
```
Authorization: ApiKey <polka_api_key>
```

**Request Body:**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Response:** `204 No Content`

**Notes:**
- Only processes "user.upgraded" events
- Upgrades user to Chirpy Red status

## Admin Endpoints

### Get Metrics
View application metrics (HTML page).

**Endpoint:** `GET /admin/metrics`

**Response:** `200 OK`
```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
</html>
```

### Reset Database
Reset users table and metrics (development only).

**Endpoint:** `POST /admin/reset`

**Response:** `200 OK`
```
users table has been cleared and hits have been reset
```

**Notes:**
- Only available when `PLATFORM` environment variable is set to "dev"

## Health Check

### Health Check
Check if the API is running.

**Endpoint:** `GET /api/healthz`

**Response:** `200 OK`
```json
{
  "status": "ok"
}
```

## Error Responses

All error responses follow this format:
```json
{
  "error": "Description of the error"
}
```

Common HTTP status codes:
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Environment Variables

Required environment variables:
- `DB_URL`: PostgreSQL database connection string
- `JWT_SECRET`: Secret key for JWT token signing
- `POLKA_KEY`: API key for Polka webhook authentication
- `PLATFORM`: Set to "dev" to enable reset endpoint

## Static Files

The application also serves static files:
- `/app/`: Main application files
- `/app/assets/`: Static assets (CSS, JS, images)

These endpoints increment the fileserver hit counter displayed in the metrics.