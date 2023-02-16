# Go Exchange

## Database [Postgres]
- [x] #1 - Design DB schema and generate SQL code with [dbdiagram](https://dbdiagram.io/)
- [x] #2 - Setup Docker Compose + Postgres + TablePlus to create DB schema
- [x] #3 - Write and Run database migration with [Migrate](https://github.com/golang-migrate/migrate)
- [x] #4 - Generate CRUD code from SQL with [SQLC](https://sqlc.dev/)
- [x] #5.1 - Write unit tests for database CRUD with random data, using [PQ](https://github.com/lib/pq) and [Testify](https://github.com/stretchr/testify)
- [x] #5.2 - Load config from file and environment variables with [Viper](https://github.com/spf13/viper)
- [x] #5.3 - Securely store passwords using Hash password with Bcrypt
- [ ] Implement database transaction
- [ ] DB transaction lock and handle deadlock
- [ ] Avoid deadlock in DB transaction. Queries order matters
- [ ] Setup Github Actions for Golang + Postgres to run automated tests

## Building RESTful HTTP JSON API [Gin]
- [ ] Implement RESTful HTTP API using [Gin](https://github.com/gin-gonic/gin)
- [ ] Mock DB for testing HTTP API and achieve 100% coverage with [GoMock](https://github.com/golang/mock)
- [ ] Implement transfer money API with a custom params validator
- [ ] Handle DB errors correctly
- [ ] Write stronger unit tests with a custom [GoMock](https://github.com/golang/mock) matcher
- [ ] Create and verify [JWT](https://github.com/golang-jwt/jwt) and [PASETO](https://github.com/o1egl/paseto) token with [UUID](https://github.com/google/uuid)
- [ ] Implement login user API that returns [JWT](https://github.com/golang-jwt/jwt) or [PASETO](https://github.com/o1egl/paseto) access token
- [ ] Implement authentication middleware and authorization rules using Gin

## Sessions, Documentation and gRPC
- [ ] Manage user session with refresh token
- [ ] Generate DB documentation page and schema SQL dump from [DB Docs](https://dbdocs.io/docs) and [DBML](https://www.dbml.org/cli/#installation)
- [ ] Introduction to gRPC
- [ ] Define gRPC API and generate Go code with [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/)
- [ ] Run a golang gRPC server and call its API
- [ ] Implement gRPC API to create and login users
- [ ] Write code once, serve both gRPC and HTTP requests with [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [ ] Extract info from gRPC metadata
- [ ] Automatic generate and serve Swagger docs from Go server with [Swagger UI](https://github.com/swagger-api/swagger-ui)
- [ ] Embed static frontend files inside Golang backend server's binary with [Statik](https://github.com/rakyll/statik)
- [ ] Validate gRPC parameters and send human/machine friendly response
- [ ] Run DB migrations directly inside Golang code
- [ ] Partial update DB record with SQLC nullable parameters
- [ ] Build gRPC update API with optional parameters
- [ ] Add authorization to protect gRPC API

---

Roadmap based on Tech School course