# 🐦 Chirpy

A Twitter-like social media platform built with Go! This project is a Twitter clone that allows users to post short messages called "chirps", manage user accounts, and interact with a RESTful API.

## 📚 About This Project

This project was built following the **Boot.Dev HTTP Servers Course**! It's a comprehensive implementation of a social media backend that demonstrates modern web development practices with Go, including:

- RESTful API design
- JWT authentication
- Database integration with PostgreSQL
- Webhook handling
- Static file serving
- Middleware implementation

## ✨ Features

- **User Management**: Create accounts, login, and update profile information
- **Chirps**: Post, read, and delete short messages (Twitter-like posts)
- **Authentication**: Secure JWT-based authentication system
- **Chirpy Red**: Premium user upgrades via webhook integration
- **Profanity Filter**: Automatic filtering of inappropriate content
- **Admin Dashboard**: Metrics and database management tools
- **Static File Serving**: Frontend asset delivery

## 🔧 Technology Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **HTTP Router**: Go standard library `http.ServeMux`
- **Database Query Builder**: [SQLC](https://sqlc.dev/)
- **Database Migrations**: [Goose](https://github.com/pressly/goose)
- **Environment Management**: [godotenv](https://github.com/joho/godotenv)
- **Password Hashing**: bcrypt
- **UUID Generation**: Google UUID library

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- [Goose](https://github.com/pressly/goose) (for database migrations)
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/05blue04/chirpy.git
   cd chirpy
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Install Goose for migrations**
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

4. **Set up environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   DB_URL=postgres://username:password@localhost/chirpy?sslmode=disable
   JWT_SECRET=your-super-secret-jwt-key
   POLKA_KEY=your-polka-api-key
   PLATFORM=dev
   ```

4. **Set up the database**
   
   Make sure PostgreSQL is running and create your database schema using your preferred migration tool.

5. **Run the application**
   ```bash
   go run .
   ```

The server will start on `http://localhost:8080`

## 📖 API Documentation

For complete API documentation including all endpoints, request/response formats, and authentication details, see:

**[API Documentation](./API_DOCS.md)**

The API documentation covers:
- User registration and authentication
- Chirp creation and management
- Token refresh and revocation
- Webhook endpoints
- Admin functionality

## 🔐 Authentication

Chirpy uses JWT (JSON Web Tokens) for authentication:

1. **Register** or **login** to receive an access token and refresh token
2. Include the access token in the `Authorization` header: `Bearer <token>`
3. Use the refresh token to get new access tokens when they expire
4. Revoke refresh tokens when logging out

## 🎯 Key Features Explained

### Chirps
- Maximum 140 characters (just like early Twitter!)
- Built-in profanity filter for words like "kerfuffle", "sharbert", and "fornax"
- Users can only delete their own chirps

### Chirpy Red
- Premium user status upgrades
- Integrated with Polka payment system via webhooks
- Automatic user status updates when payments are processed

### Admin Features
- Metrics dashboard showing application usage
- Database reset functionality (development mode only)
- Request monitoring and hit counting

## 🧪 Development

### Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migrations. Migration files are located in `sql/schema/`.

**Common migration commands:**
```bash
# Apply all pending migrations
goose -dir sql/schema postgres $DB_URL up

# Rollback the last migration
goose -dir sql/schema postgres $DB_URL down

# Check migration status
goose -dir sql/schema postgres $DB_URL status

# Create a new migration
goose -dir sql/schema create migration_name sql
```

### Running in Development Mode

Set `PLATFORM=dev` in your `.env` file to enable:
- Database reset endpoint (`POST /admin/reset`)
- Additional logging and debugging features

### Testing

```bash
go test ./...
```
## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- **[Boot.Dev](https://boot.dev)** for the excellent HTTP Servers course that guided this project
- The Go community for amazing tools and libraries
- PostgreSQL team for the robust database system

*Built with ❤️ following Boot.Dev's HTTP Servers course
