# Social Media API

REST API for a social media platform built with Go, Fiber, and MongoDB.

## Tech Stack

- **Go** + **Fiber v3** — HTTP framework
- **MongoDB** — database
- **Cloudinary** — media storage
- **JWT** — authentication

## Project Structure
```
internal/
├── handler/        # HTTP handlers
├── service/        # Business logic
├── repository/     # Database layer
├── middleware/     # Auth, etc.
├── model/          # MongoDB models
├── dto/            # Request/Response structs
└── helpers/        # Utilities
```

## Features

- User authentication (JWT)
- Follow/unfollow users
- User suggestions based on social graph
- Create posts with images
- Profile management

## Getting Started

### Prerequisites
- Go 1.21+
- MongoDB
- Cloudinary account

### Environment Variables
```
PORT=8080
MONGO_URL=""
MONGO_DB=""
JWT_SECRET=""
JWT_ISSUER=""
CLD_API_KEY=""
CLD_API_SECRET=""
CLD_CLOUD_NAME=""
```
### Run

git clone https://github.com/tu-usuario/tu-repo
cd tu-repo
go mod tidy
go run cmd/main.go

## License
MIT © [Arnulfo Vargas Mejia](https://github.com/ArnulfoVargas)
