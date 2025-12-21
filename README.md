# gopher-social
 Gopher Social is a

A social media REST API built with Go

---

## Tech Stack
- Go (Golang, Fiber)
- PostgreSQL
- Redis
- Docker & Docker Compose
- Other (GORM, ginkgo)

---
## CMI

### Run Docker compose (database, redis)
```
docker-compose up --build -d
```

```
docker-compose down
```
### Run the API Server
```go
go run ./cmd/api
```

### Run Database Migrations
```go
go run ./cmd/api/ migrate up
```

```go
go run ./cmd/api/ migrate down
```

### Run Test
```go
ginkgo ./cmd/api/
```