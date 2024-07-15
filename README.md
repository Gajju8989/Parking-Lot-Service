# Parking Lot Service 

Author: Gajendra Singh

## Description
A HTTP service to manage parking lots, allowing parking of various vehicle types (Motorcycles, Cars, Buses)
with different spot sizes and tariffs. Built using Echo and GORM, 
with PostgreSQL as the database. Features real-time free slot tracking
and generation of parking tickets and receipts, including computation of 
parking fees based on the lot-specific tariff models.

## Dependencies
- Go 1.21.6
     - github.com/joho/godotenv v1.5.1
     - github.com/labstack/echo/v4 v4.12.0
     - github.com/swaggo/echo-swagger v1.4.1
     - github.com/swaggo/swag v1.16.3
     - gorm.io/driver/postgres v1.5.9
     - gorm.io/gorm v1.25.11
     - (Other  dependencies listed in `go.mod` file)
## Project Structure
The project is structured as follows:
 ```graphql
parking_lot_service/
├── docs/
│ ├── docs.go # Documentation generation script
│ ├── swagger.json # Swagger JSON file
│ └── swagger.yaml # Swagger YAML file
├── internal/
│ ├── database/
│ │ └── postgresql/
│ │ ├── config/
│ │ │ └── postgres_config.go # PostgreSQL configuration
│ │ └── migration/
│ │ └── migration.go # Database migration script
│ ├── di/
│ │ └── container.go # Dependency Injection container setup
│ ├── genericresponse/
│ │ └── genericresponse.go # Generic HTTP response handling
│ ├── handler/
│ │ ├── handler.go # HTTP handler definitions
│ │ ├── handler_get_parking_space_impl.go # Implementation of Get Parking Space handler
│ │ ├── handler_park_vehicle_impl.go # Implementation of Park Vehicle handler
│ │ └── handler_un_park_vehicle_impl.go # Implementation of Unpark Vehicle handler
│ ├── repo/
│ │ ├── models/
│ │ │ └── models.go # Data models
│ │ ├── repo.go # Repository interface definitions
│ │ └── repo_impl.go # Repository implementations
│ ├── router/
│ │ ├── router.go # HTTP router setup
│ │ └── router_impl.go # HTTP router implementations
│ └── service/
│ ├── model/
│ │ ├── model.go # Service models
│ │ └── commons.go # Common utilities for services
│ ├── service.go # Service interface definitions
│ ├── service_get_parking_space_impl.go # Implementation of Get Parking Space service
│ ├── service_un_park_vehicle_impl.go # Implementation of Unpark Vehicle service
│ └── service_un_park_vehicle_impl_test.go # Unit tests for Unpark Vehicle service
├── go.mod # Go module file
├── local.env # Environment variables file
├── main.go # Main application entry point
└── README.md # Project documentation (you are here)
```



## Installation and Usage
### Install Dependencies:
   ```bash
    go mod download
   ```


### Local Environment Setup
Before running the application, you need to set up your local environment variables by 
creating a .env file in the project root directory. Add the following configuration to the .env file:
```text
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=123
DB_NAME=parking_lot_service
```

Replace the values with your database configuration.

### Run Server
  ```bash
go run main.go 
```
## Testing
To run tests of service:
```bash
 go test -v ./internal/service
```

## API Documentation

### Postman API Documentation
Here's the link to the API documentation for HTTP endpoints:
- [Postman API Documentation](https://documenter.getpostman.com/view/29203481/2sA3kPqQYw)

### Swagger API Documentation
To view the Swagger documentation:
1.Start the server using 
```bash 
 go run main.go
```
2.Navigate to Swagger UI after starting the server to explore and interact with the API endpoints.