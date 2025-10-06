# Multi-Datasource Go

A Go application demonstrating connectivity and operations with multiple database systems (MySQL, PostgreSQL, and Oracle XE) using a clean architecture pattern.

## Features

- 🔌 Multiple database connections with connection pooling
- 🌐 RESTful API using Gin framework
- 🎯 Domain-driven design with repository pattern
- 📝 YAML-based configuration with Viper
- 🐳 Docker Compose setup for local development
- ⚡ Built-in health checks
- 🏗️ Clean architecture (domain, repository, service layers)
- **Auto Table Creation**: Automatically ensures database tables exist on startup
- **Modular Design**: Each datasource is independently managed

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git

## Project Structure

```
multi-datasource-go/
├─ cmd/
│  └─ api/
│     └─ main.go           # Application entry point
├─ internal/
│  ├─ config/
│  │  └─ config.go         # Configuration management
│  ├─ db/
│  │  ├─ mysql.go          # MySQL connection
│  │  ├─ postgres.go       # PostgreSQL connection
│  │  └─ oracle.go         # Oracle connection
│  ├─ http/
│  │  └─ handlers.go       # Gin routes + handlers
│  ├─ domain/
│  │  ├─ model.go          # User, Company, Brand structs
│  │  ├─ repo.go           # UserRepo, CompanyRepo, BrandRepo interfaces
│  │  └─ service.go        # UserService, CompanyService, BrandService
│  └─ repo/
│     ├─ mysql_user_repo.go    # MySQLUserRepo (users)
│     ├─ pg_company_repo.go    # PGCompanyRepo (companies)
│     └─ oracle_brand_repo.go  # OracleBrandRepo (brands)
├─ application.yaml            # Application configuration
├─ go.mod                      # Go module dependencies
├─ go.sum
├─ docker-compose-multiple-db.yml  
└─ README.md

```

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/HenryXiloj/demos-go.git
cd demos-go/multi-datasource-go
```

### 2. Initialize Go Module and Install Dependencies

```bash
# Initialize module
go mod init multi-datasource-go

# Install required packages
go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/go-sql-driver/mysql
go get github.com/sijms/go-ora/v2@latest
go get github.com/spf13/viper
```

### 3. Configuration Setup

Create an `application.yaml` file:

```yaml
app:
  httpPort: 9000
  requestTimeoutSec: 5

mysql:
  enabled: true
  dsn: "test:test_pass@tcp(127.0.0.1:3306)/test_db?parseTime=true"
  maxOpenConns: 50
  maxIdleConns: 10
  connMaxLifetimeMin: 30
  connMaxIdleMin: 5

postgres:
  enabled: true
  dsn: "postgres://postgre_test:postgre_test@127.0.0.1:5432/postgre_test?sslmode=disable"
  maxOpenConns: 50
  maxIdleConns: 10
  connMaxLifetimeMin: 30
  connMaxIdleMin: 5

oracle:
  enabled: true
  dsn: "oracle://app:apppw@127.0.0.1:1521/xepdb1"
  maxOpenConns: 30
  maxIdleConns: 5
  connMaxLifetimeMin: 30
  connMaxIdleMin: 5
```

### 4. Start Database Services

```bash
docker compose -f docker-compose-multiple-db.yml up -d
```

Wait for all containers to be healthy:

```bash
docker-compose ps
```

### 5. Run the Application

```bash
go mod tidy
go run cmd/api/main.go
```

The API will be available at `http://localhost:9000`

You should see output similar to:

```
2025/10/05 18:26:29 ✅ ensured MySQL table: users
2025/10/05 18:26:29 ✅ ensured Postgres table: companies
2025/10/05 18:26:29 ✅ ensured Oracle table: brands
[GIN-debug] POST   /api/v1/users             --> multi-datasource-go/internal/http.(*Handlers).createUser-fm (2 handlers)
[GIN-debug] POST   /api/v2/companies         --> multi-datasource-go/internal/http.(*Handlers).createCompany-fm (2 handlers)
[GIN-debug] POST   /api/v3/brands            --> multi-datasource-go/internal/http.(*Handlers).createBrand-fm (2 handlers)
2025/10/05 18:26:29 listening on :9000
```

## Quick Health Checks

Test database connectivity without installing local clients:

### MySQL Health Check

```bash
# Ping MySQL and check version
docker exec -it mysql84 bash -lc "mysql -uroot -p$MYSQL_ROOT_PASSWORD -e 'SELECT VERSION();'"
```

Expected output:
```
+-----------+
| VERSION() |
+-----------+
| 8.4.x     |
+-----------+
```

### PostgreSQL Health Check

```bash
# Check PostgreSQL version and connectivity
docker exec -it postgre_test bash -lc "psql -U postgre_test -d postgre_test -c 'SELECT version();'"
```

Expected output:
```
                                                version                                                
--------------------------------------------------------------------------------------------------------
 PostgreSQL 15.x on x86_64-pc-linux-gnu, compiled by gcc...
```

### Oracle XE Health Check

```bash
# Test Oracle connection and basic query
docker exec -it oracle-xe bash -lc "echo 'SELECT 1 FROM dual;' | sqlplus -s app/apppw@localhost:1521/xepdb1"
```

Expected output:
```
         1
----------
         1
```

## 📡 API Endpoints

### Users API (MySQL)

**Create User**
```bash
curl -X POST http://localhost:9000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Henry","lastName":"x"}'
```

Response:
```json
{"id":1}
```

### Companies API (PostgreSQL)

**Create Company**
```bash
curl -X POST http://localhost:9000/api/v2/companies \
  -H "Content-Type: application/json" \
  -d '{"name":"test"}'
```

Response:
```json
{"id":1}
```

### Brands API (Oracle)

**Create Brand**
```bash
curl -X POST http://localhost:9000/api/v3/brands \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme"}'
```

Response:
```json
{"id":1}
```

## 🧪 Testing

Run the application and test each endpoint:

```bash
# Terminal 1: Start the server
go run cmd/api/main.go

# Terminal 2: Test endpoints
# Create a user (MySQL)
curl -X POST http://localhost:9000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Henry","lastName":"x"}'

# Create a company (PostgreSQL)
curl -X POST http://localhost:9000/api/v2/companies \
  -H "Content-Type: application/json" \
  -d '{"name":"test"}'

# Create a brand (Oracle)
curl -X POST http://localhost:9000/api/v3/brands \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme"}'
```

## 🏗️ Architecture

This project follows a **layered architecture** pattern:

### 1. **Handlers Layer** (`internal/http`)
- Handles HTTP requests and responses
- Validates input data
- Delegates business logic to services

### 2. **Service Layer** (`internal/service`)
- Contains business logic
- Orchestrates operations between repositories
- Independent of HTTP concerns

### 3. **Repository Layer** (`internal/repository`)
- Handles database operations
- Abstracts database-specific implementations
- One repository per entity/database

### 4. **Models Layer** (`internal/models`)
- Defines data structures
- Represents database entities

## 🗄️ Database Schema

### MySQL - Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### PostgreSQL - Companies Table
```sql
CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Oracle - Brands Table
```sql
CREATE TABLE brands (
    id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR2(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Technologies Used

- **[Go](https://golang.org/)** - Programming language
- **[Gin](https://github.com/gin-gonic/gin)** - HTTP web framework
- **[Viper](https://github.com/spf13/viper)** - Configuration management
- **[MySQL Driver](https://github.com/go-sql-driver/mysql)** - MySQL database driver
- **[pgx](https://github.com/jackc/pgx)** - PostgreSQL driver and toolkit
- **[go-ora](https://github.com/sijms/go-ora)** - Oracle database driver

## 📝 Notes

- Tables are automatically created on application startup if they don't exist
- Each database connection is managed independently
- The application uses connection pooling for optimal performance
- Each database can be enabled/disabled via configuration

## 🙏 Acknowledgments

- Clean Architecture principles by Robert C. Martin
- Go community for excellent database drivers

- Gin framework for the lightweight HTTP router

