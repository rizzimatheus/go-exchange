# Go Exchange

## 1. Database [Postgres]
- [x] v1.1.0 - Design DB schema and generate SQL code with [dbdiagram](https://dbdiagram.io/)
- [x] v1.2.0 - Setup Docker Compose + Postgres + TablePlus to create DB schema
- [x] v1.3.0 - Write and Run database migration with [Migrate](https://github.com/golang-migrate/migrate)
- [x] v1.4.0 - Generate CRUD code from SQL with [SQLC](https://sqlc.dev/)
- [x] v1.5.1 - Write unit tests for database CRUD with random data, using [PQ](https://github.com/lib/pq) and [Testify](https://github.com/stretchr/testify)
- [x] v1.5.2 - Load config from file and environment variables with [Viper](https://github.com/spf13/viper)
- [x] v1.5.3 - Securely store passwords using Hash password with Bcrypt
- [x] v1.6.0 - Implement database transaction
- [x] v1.7.0 - Setup Github Actions for Golang + Postgres to run automated tests

## 2. Building RESTful HTTP JSON API [Gin]
- [x] v2.1.0 - Implement RESTful HTTP API using [Gin](https://github.com/gin-gonic/gin) and custom params validator
- [x] v2.2.1 - Run DB migrations directly inside Golang code
- [x] v2.2.2 - Build a minimal Golang Docker image with a multistage Dockerfile
- [x] v2.3.1 - Mock DB for testing HTTP API and achieve 100% coverage with [GoMock](https://github.com/golang/mock)
- [x] v2.3.2 - Write stronger unit tests with a custom [GoMock](https://github.com/golang/mock) matcher
- [x] v2.4.0 - Create and verify [JWT](https://github.com/golang-jwt/jwt) and [PASETO](https://github.com/o1egl/paseto) token with [UUID](https://github.com/google/uuid)
- [x] v2.5.0 - Implement login user API that returns [JWT](https://github.com/golang-jwt/jwt) or [PASETO](https://github.com/o1egl/paseto) access token
- [x] v2.6.0 - Implement authentication middleware and authorization rules using [Gin](https://github.com/gin-gonic/gin)
- [x] v2.7.0 - Manage user session with refresh token

## 3. gRPC
- [x] v3.1.0 - Define gRPC API and generate Go code with [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/)
- [x] v3.2.0 - Run a golang gRPC server and call its API
- [x] v3.3.0 - Implement gRPC API to create, login and update users
- [x] v3.4.0 - Extract info from gRPC metadata
- [x] v3.5.0 - Add authorization to protect gRPC API
- [x] v3.6.0 - Validate gRPC parameters and send human/machine friendly response
- [x] v3.7.0 - Write code once, serve both gRPC and HTTP requests with [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [x] v3.8.0 - Write structured logs for gRPC APIs with [ZeroLog](https://github.com/rs/zerolog)
- [x] v3.9.0 - Write HTTP logger middleware with [ZeroLog](https://github.com/rs/zerolog)

## 4. Documentation
- [x] v4.1.0 - Generate DB documentation page and schema SQL dump from [DB Docs](https://dbdocs.io/docs) and [DBML](https://www.dbml.org/cli/#installation)
- [x] v4.2.0 - Automatic generate and serve Swagger docs from Go server with [Swagger UI](https://github.com/swagger-api/swagger-ui)
- [x] v4.3.0 - Embed static frontend files inside Golang backend server's binary with [Statik](https://github.com/rakyll/statik)

---

Roadmap based on Tech School course