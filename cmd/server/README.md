# API

## Features

- **Configuration**: Support for local, development, and QA environments.
- **API Authentication**: Secure access to api based on rbac.
- **Database Migrations**: Automated database schema management.
- **OpenAPI Documentation**: RESTful API documentation generated using Swagger.

### Module layer

- **Middleware**: Custom middleware for logging and authentication.
- **Handler Layer**: Handler request, validator request body.
- **Service Layer**: Separation of concerns with a clean architecture design.
- **Data Layer**: SQL and redis control data.
- **API Layer**: Third party api.
- **Data Models Layer**: Defined data models api response.
- **Error Handling**: Centralized error handling and response formatting.

### Technologies Used

- **Go (Golang)**: Primary programming language.
- **Fx**: Application dependency injection and running entry.
- **Fiber**: Web framework for building RESTful APIs.
- **Ent ORM**: For database schema modeling.
- **Atlas**: For database migration management.
- **Swagger (OpenAPI)**: For API documentation.
- **MySQL**: Relational database for storing data.
- **Casbin**: For role-based access control and authorization.
- **Redis**: For caching and session management.
- **MQTT**: For communication with devices.
- **Cos (Tencent Cloud)**: For cloud storage integration.
- **OpenTelemetry**: For distributed tracing and metrics collection.
- **Go-Redis**: For Redis client integration.
- **Go-Validator**: For request validation.
- **Zap**: For structured logging.
- **Lumberjack**: For log rotation and management.
- **Docker**: For containerization testing,developing and deployment.
- **Testcontainers**: For containerized testing.
- **Testify**: For unit testing, assert.

## Documents

- **Swagger API Docs**: Generated using Swagger.
- **OpenAPI Specification**: See `docs/swagger.json` for the full API definition.

```bash
# Generate OpenAPI documentation
# API service router update should be update openapi json
swag init -d ./cmd/api -g router.go --outputTypes json --pdl 3
```

## Develope project

### Prerequisites

- Go 1.24 or higher
- MySQL 8
- Redis >= 7
- MQTT Broker
- Docker (for local development and unit test)
- OpenTelemetry Collector (optional, for distributed tracing)

### Run the Application

Application should be run in configuration, description in `configs/api_config_template.yaml`

```bash
go run cmd/api/main.go -c ./devops/configs/api/config_local.yaml
```
