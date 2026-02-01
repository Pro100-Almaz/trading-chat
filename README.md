# Trading Chat Backend

A RESTful API backend built with Go following Clean Architecture principles. Features user authentication via email/password and Google OAuth 2.0, JWT-based authorization, and user management.

## Features

- **User Registration & Login** - Traditional email/password authentication
- **Google OAuth 2.0** - Social login integration
- **JWT Authentication** - Access and refresh token system
- **User CRUD Operations** - Full user management API
- **Clean Architecture** - Layered, testable, maintainable codebase
- **PostgreSQL Database** - Persistent data storage
- **Swagger Documentation** - Interactive API documentation

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.19+ |
| Router | gorilla/mux |
| Database | PostgreSQL with sqlx |
| Auth | JWT (golang-jwt), bcrypt, OAuth2 |
| Config | Viper |
| Logging | Logrus |
| API Docs | Swagger (swaggo) |

## Project Structure

```
backend/
├── cmd/
│   └── main.go                 # Application entry point
├── api/
│   ├── controller/             # HTTP request handlers
│   │   ├── signup.go
│   │   ├── login.go
│   │   ├── google.go
│   │   ├── refresh_token.go
│   │   ├── user.go
│   │   └── error_response.go
│   ├── middleware/             # HTTP middleware
│   │   ├── jwt_auth.go         # JWT validation middleware
│   │   └── logger.go           # Request logging
│   └── route/                  # Route definitions
│       ├── route.go
│       ├── signup.go
│       ├── login.go
│       ├── google.go
│       ├── refresh_token.go
│       └── user.go
├── bootstrap/                  # App initialization
│   ├── app.go                  # Application factory
│   ├── env.go                  # Environment config
│   └── database.go             # Database connection
├── domain/                     # Business models & interfaces
│   ├── user.go
│   ├── signup.go
│   ├── login.go
│   ├── refresh_token.go
│   ├── google.go
│   ├── jwt_custom.go
│   └── error_response.go
├── repository/                 # Data access layer
│   └── user_repository.go
├── usecase/                    # Business logic layer
│   ├── signup_usecase.go
│   ├── login_usecase.go
│   ├── google_usecase.go
│   ├── refresh_token_usecase.go
│   └── user_usecase.go
├── internal/
│   └── tokenutil/              # JWT utilities
│       └── tokenutil.go
├── utils/
│   └── util.go                 # Helper functions
├── .env.example                # Environment template
├── Dockerfile                  # Docker configuration
├── go.mod
└── go.sum
```

## Architecture

The project follows Clean Architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Layer (api/)                        │
│         Controllers → Routes → Middleware                    │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                  Business Logic (usecase/)                   │
│    SignupUseCase, LoginUseCase, UserUseCase, GoogleUseCase  │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                Data Access Layer (repository/)               │
│                      UserRepository                          │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                    PostgreSQL Database                       │
└─────────────────────────────────────────────────────────────┘
```

**Layer Responsibilities:**
- **Controllers**: Handle HTTP requests/responses, input validation
- **Use Cases**: Business logic, orchestrate operations
- **Repositories**: Database queries, data persistence
- **Domain**: Data models, interfaces, error definitions

## API Documentation (Swagger)

Interactive API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

To regenerate Swagger docs after making changes to API annotations:

```bash
# Install swag CLI (one-time)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/main.go -o docs
```

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/signup` | Register new user |
| `POST` | `/api/login` | Login with email/password |
| `GET` | `/api/google/login` | Initiate Google OAuth flow |
| `GET` | `/api/google/callback` | Google OAuth callback |
| `POST` | `/api/refresh_token` | Refresh JWT tokens |

### Protected Endpoints (Require JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/user/all` | Get all users |
| `GET` | `/api/user` | Get current user profile |
| `PUT` | `/api/user` | Update current user |
| `DELETE` | `/api/user` | Delete current user |

### Request/Response Examples

#### Sign Up
```bash
POST /api/signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```
Response:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Login
```bash
POST /api/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securePassword123"
}
```
Response:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Get User Profile (Protected)
```bash
GET /api/user
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```
Response:
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "",
  "profile_picture": "",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### Update User (Protected)
```bash
PUT /api/user
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "name": "John Updated",
  "phone": "+1234567890"
}
```

#### Refresh Token
```bash
POST /api/refresh_token
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```
Response:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

## Getting Started

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 12+
- (Optional) Google Cloud Console project for OAuth

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your configuration (see Configuration section below).

4. **Set up PostgreSQL database**
   ```bash
   # Connect to PostgreSQL
   psql -U postgres

   # Create database
   CREATE DATABASE trading_chat;
   ```
   The application will auto-migrate the `users` table on startup.

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

   The server will start at `http://localhost:8080` (or your configured port).

### Configuration

Create a `.env` file in the project root with the following variables:

```env
# Application
APP_ENV=development
SERVER_ADDRESS=:8080
PORT=8080
CONTEXT_TIMEOUT=2

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your_password
DB_NAME=trading_chat

# JWT Configuration
ACCESS_TOKEN_EXPIRY_HOUR=2
REFRESH_TOKEN_EXPIRY_HOUR=168
ACCESS_TOKEN_SECRET=your_access_token_secret_here
REFRESH_TOKEN_SECRET=your_refresh_token_secret_here

# Google OAuth (Optional)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment mode (`development`/`production`) | - |
| `SERVER_ADDRESS` | Server listen address | `:8080` |
| `PORT` | Server port | `8080` |
| `CONTEXT_TIMEOUT` | Request timeout in seconds | `2` |
| `DB_HOST` | PostgreSQL host | - |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL username | - |
| `DB_PASS` | PostgreSQL password | - |
| `DB_NAME` | Database name | - |
| `ACCESS_TOKEN_EXPIRY_HOUR` | Access token lifetime (hours) | `2` |
| `REFRESH_TOKEN_EXPIRY_HOUR` | Refresh token lifetime (hours) | `168` (7 days) |
| `ACCESS_TOKEN_SECRET` | JWT signing key for access tokens | - |
| `REFRESH_TOKEN_SECRET` | JWT signing key for refresh tokens | - |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | - |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | - |

### Google OAuth Setup (Optional)

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Navigate to **APIs & Services** > **Credentials**
4. Click **Create Credentials** > **OAuth 2.0 Client ID**
5. Set application type to **Web application**
6. Add authorized redirect URI: `http://localhost:8080/api/google/callback`
7. Copy Client ID and Client Secret to your `.env` file

## Running with Docker

### Using Docker Compose (Recommended)

The easiest way to run the application with PostgreSQL:

```bash
# Start PostgreSQL database
docker-compose up -d

# Check if database is ready
docker-compose ps

# View database logs
docker-compose logs db
```

This starts a PostgreSQL 16 container with:
- **Database**: `myapp`
- **User**: `myapp`
- **Password**: `myapp_pass`
- **Port**: `5433` (mapped to container's `5432`)

Update your `.env` file to match:
```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=myapp
DB_PASS=myapp_pass
DB_NAME=myapp
```

Then run the application:
```bash
go run cmd/main.go
```

To stop the database:
```bash
docker-compose down

# To also remove the data volume:
docker-compose down -v
```

### Using Docker Only

```bash
# Build the image
docker build -t trading-chat-backend .

# Run the container
docker run -p 8080:8080 --env-file .env trading-chat-backend
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./usecase/...
```

## Database Schema

The application automatically creates the following table:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    google_id VARCHAR(255) DEFAULT '',
    profile_picture VARCHAR(255) DEFAULT '',
    name VARCHAR(255) DEFAULT '',
    password VARCHAR(255) DEFAULT '',
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
```

## Authentication Flow

### Email/Password Flow
1. User registers via `POST /api/signup` with name, email, password
2. Password is hashed with bcrypt and stored
3. Server returns access token (2h) and refresh token (7d)
4. Client includes access token in `Authorization: Bearer <token>` header
5. When access token expires, use refresh token to get new tokens

### Google OAuth Flow
1. Client redirects to `GET /api/google/login`
2. User authenticates with Google
3. Google redirects to `/api/google/callback`
4. Server creates/updates user and sets auth cookies
5. Client is redirected to `/profile` with authenticated session

## Error Responses

All errors return JSON in this format:
```json
{
  "message": "Error description"
}
```

Common HTTP status codes:
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (invalid/missing token)
- `404` - Not Found (resource doesn't exist)
- `500` - Internal Server Error

## License

This project is open source and available under the MIT License.
