# Gopher Social

Gopher Social is a RESTful API for a social media platform, built with Go. It provides endpoints for user authentication, posting content, following users, and more.

---

## Features
- User authentication and role-based access control
- Create, read, update, and delete posts
- Follow/unfollow users
- Commenting and tagging on posts
- Rate limiting for API requests
- Dockerized setup for easy deployment

---

## Tech Stack
- **Go** (Golang, Fiber)
- **PostgreSQL**
- **Redis**
- **Docker & Docker Compose**
- Other: GORM, Ginkgo

---

## API Documentation
The API is documented using Swagger. You can view the documentation by running the server and navigating to:
```
http://localhost:3000/swagger/index.html
```

---

## Setup Instructions

### Prerequisites
- Install [Go](https://golang.org/doc/install)
- Install [Docker](https://www.docker.com/get-started)
- Install [Ginkgo](https://onsi.github.io/ginkgo/)

### Run Docker Compose (Database, Redis)
To start the required services:
```
docker-compose up --build -d
```
To stop the services:
```
docker-compose down
```

### Run the API Server
Start the API server with:
```
go run ./cmd/api
```

### Run Database Migrations
To apply migrations:
```
go run ./cmd/api migrate up
```
To roll back migrations:
```
go run ./cmd/api migrate down
```

---

## Testing
Run the test suite using Ginkgo:
```
ginkgo ./cmd/api/
```
Ensure all tests pass before committing changes.

---

## Contributing
We welcome contributions! To contribute:
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes with clear messages.
4. Open a pull request.

---