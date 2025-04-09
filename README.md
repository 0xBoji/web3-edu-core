# Web3 Education Platform API

A Golang API for an educational platform with courses, lessons, users, enrollments, and progress tracking.

## Features

- User authentication and authorization
- Course management
- Lesson management
- User enrollment and progress tracking
- Internationalization (i18n) support
- Redis caching
- Swagger API documentation

## Tech Stack

- Go 1.21+
- Gin Web Framework
- GORM (PostgreSQL)
- Redis
- JWT Authentication
- Swagger
- Docker

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go            # Entry point
├── config/
│   ├── app.ini                # Configuration file
│   └── config.go              # Configuration loader
├── docs/                      # Swagger documentation
├── internal/
│   ├── api/
│   │   ├── docs.go            # Swagger docs
│   │   ├── middleware/        # Middleware components
│   │   ├── server.go          # Server setup
│   │   └── v1/                # API version 1
│   │       ├── handlers/      # Request handlers
│   │       └── routes/        # API routes
│   ├── database/
│   │   ├── postgres/          # PostgreSQL connection
│   │   └── redis/             # Redis connection and cache
│   ├── domain/
│   │   ├── models/            # Database models
│   │   ├── repositories/      # Data access layer
│   │   └── services/          # Business logic
│   └── utils/                 # Utility functions
├── locales/                   # Translation files
├── migrations/                # SQL migration files
├── scripts/                   # Helper scripts
├── .dockerignore
├── .env                       # Environment variables
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Running with Docker

1. Clone the repository:

```bash
git clone https://github.com/0xBoji/web3-edu-core.git
cd web3-edu-core
```

2. Start the application:

```bash
docker-compose up -d
```

3. Access the API at http://localhost:8003

4. Access the Swagger documentation at http://localhost:8003/swagger/index.html

5. Access PgAdmin at http://localhost:5050 (Email: admin@web3edu.com, Password: admin)

### Development

1. Install Go 1.21 or later
2. Install PostgreSQL and Redis
3. Clone the repository
4. Install dependencies:

```bash
make deps
```

5. Run the application:

```bash
make run
```

### Using the Makefile

The project includes a Makefile to simplify common development tasks:

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linters
make lint

# Generate Swagger documentation
make swagger

# Build Docker image
make docker-build

# Start all services with Docker Compose
make docker-compose-up

# Stop all services with Docker Compose
make docker-compose-down

# Run database migrations
make migrate

# Run specific migration up
make migrate-up version=20230101000001

# Run specific migration down
make migrate-down version=20230101000001

# Download dependencies
make deps

# Tidy up the go.mod file
make tidy
```

Run `make help` to see all available commands.

## API Endpoints

### Authentication
- POST   /api/v1/auth/register          - Register a new user
- POST   /api/v1/auth/login             - Login
- POST   /api/v1/auth/refresh-token     - Refresh token
- POST   /api/v1/auth/logout            - Logout
- POST   /api/v1/auth/forgot-password   - Forgot password
- POST   /api/v1/auth/reset-password    - Reset password

### User Management
- GET    /api/v1/users/me               - Get current user information
- PUT    /api/v1/users/me               - Update current user information
- PATCH  /api/v1/users/me/password      - Change password
- GET    /api/v1/users/me/enrollments   - List enrolled courses

### Categories
- GET    /api/v1/categories             - Get all categories
- GET    /api/v1/categories/{id}        - Get category details

### Courses
- GET    /api/v1/courses                - Get all courses (with filter, pagination)
- GET    /api/v1/courses/featured       - Get featured courses
- GET    /api/v1/courses/{id}           - Get course details
- GET    /api/v1/courses/{id}/lessons   - Get course lessons
- POST   /api/v1/courses/{id}/enroll    - Enroll in a course
- GET    /api/v1/courses/{id}/reviews   - Get course reviews
- POST   /api/v1/courses/{id}/reviews   - Add a review to a course

### Lessons
- GET    /api/v1/lessons/{id}            - Get lesson details
- GET    /api/v1/lessons/{id}/progress   - Get lesson progress
- POST   /api/v1/lessons/{id}/progress   - Update lesson progress
- POST   /api/v1/lessons/{id}/complete   - Mark lesson as completed

### Admin APIs
- POST   /api/v1/admin/courses           - Create a new course
- PUT    /api/v1/admin/courses/{id}      - Update a course
- DELETE /api/v1/admin/courses/{id}      - Delete a course
- POST   /api/v1/admin/lessons           - Create a new lesson
- PUT    /api/v1/admin/lessons/{id}      - Update a lesson
- DELETE /api/v1/admin/lessons/{id}      - Delete a lesson
- GET    /api/v1/admin/users             - Manage users
- PUT    /api/v1/admin/users/{id}/role   - Update user role

### Internationalization
- GET    /api/v1/i18n/{language}         - Get translations for a specific language

## License

This project is licensed under the MIT License - see the LICENSE file for details.
