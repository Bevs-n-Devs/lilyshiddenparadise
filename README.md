# Lily's Hidden Paradise

Lily's Hidden Paradise is a landlord-tenant CRM tool that allows landlords to manage their properties and tenants. The application provides secure dashboards for both landlords and tenants, with features for managing applications, messaging, and property information.

**Future iterations will also allow tenants to manage their leases and payments.**

*This tool is built using Golang, JavaScript and PostgreSQL.*

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.16+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/postgresql-12+-blue.svg)

## Table of Contents
- [Key Features](#key-features)
- [Project Structure](#project-structure)
- [Components](#components)
  - [Handlers](#handlers)
  - [Middleware](#middleware)
  - [Utils](#utils)
  - [Database](#database)
- [Test Coverage](#test-coverage)
- [Setup Instructions](#setup-instructions)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
  - [Running Tests](#running-tests)
- [Contributing](#contributing)

## Key Features

- HTML templating with `http/template` standard library
- Web authentication using sessions & CSRF protection
- Protected dashboards for landlords and tenants (using middleware)
- Creating a database connection using `database/sql` standard library and `github.com/lib/pq` for the PostgreSQL driver
- Password hashing using `golang.org/x/crypto/bcrypt`
- Encrypting & decrypting database information 
- Message platform for landlords and tenants
- Database stubbing for testing

## Project Structure

```
lilyshiddenparadise/
├── db/                 # Database connection and queries
├── env/                # Environment configuration
├── handlers/           # HTTP request handlers
│   └── backup/         # Backup handlers
├── logs/               # Logging functionality
├── middleware/         # Authentication and session middleware
├── static/             # Static assets (CSS, JS, images)
├── templates/          # HTML templates
├── testutil/           # Testing utilities
├── utils/              # Utility functions
├── main.go             # Application entry point
└── README.md           # This file
```

## Components

### Handlers

The handlers package contains HTTP request handlers that process incoming requests and return responses. Key handlers include:

- **Home**: Serves the landing page
- **Login**: Handles user authentication for landlords and tenants
- **Landlord Dashboard**: Manages landlord-specific views and actions
  - Property management
  - Tenant applications
  - Tenant management
  - Messaging
- **Tenant Dashboard**: Manages tenant-specific views and actions
  - Account management
  - Messaging
  - Password updates

Each handler follows a similar pattern:
1. Validates the request (authentication, form data)
2. Processes the request (database operations, business logic)
3. Renders the appropriate template or redirects

### Middleware

The middleware package provides authentication and session management:

- **session.go**: Contains functions for authenticating requests
  - `AuthenticateLandlordRequest`: Validates landlord session and CSRF tokens
  - `AuthenticateTenantRequest`: Validates tenant session and CSRF tokens
  - These functions check if the session token and CSRF token in the request are valid by querying the database

- **cookies.go**: Manages cookie creation and handling
  - Dashboard session cookies (e.g., `LandlordDashboardSessionCookie`, `TenantDashboardSessionCookie`)
  - CSRF token cookies (e.g., `LandlordDashboardCSRFTokenCookie`, `TenantDashboardCSRFTokenCookie`)
  - Cookie deletion for logout functionality (e.g., `DeleteLandlordSessionCookie`, `DeleteTenantCSRFCookie`)
  - Each function sets appropriate cookie parameters like path, expiry, and HttpOnly flag

### Utils

The utils package provides utility functions used throughout the application:

- **utils.go**: Core utility functions
  - `HashedPassword`: Generates bcrypt hashed passwords with a cost factor of 10
  - `CheckPasswordHash`: Verifies passwords against hashes using bcrypt's comparison
  - `HashData`: Creates SHA-256 hashes for data (used for non-password sensitive data)
  - `Encrypt`/`Decrypt`: Encrypts and decrypts sensitive data using AES-GCM
  - `GenerateToken`: Creates secure random tokens for sessions and CSRF protection
  - `ValidateAge`: Ensures users are 18+ by comparing date of birth with current date
  - Various form validation functions that check conditional requirements

- **types.go**: Defines types and initializes encryption
  - Sets up the master encryption key from environment variables
  - Provides initialization function for the encryption system

### Database

The db package handles database connections and queries:

- PostgreSQL connection management
- User authentication queries
- Property and tenant data management
- Session token validation
- Message storage and retrieval

## Test Coverage

The application includes comprehensive tests for critical components:

### Middleware Tests (26.2% coverage)
- **cookies_test.go**:
  - Tests cookie creation for different dashboard paths
  - Validates cookie properties (expiry, path, HttpOnly flag)
  - Tests cookie deletion functionality
  - Tests specialized cookies for messaging and logout

- **session_test.go**:
  - Tests authentication with missing session tokens
  - Tests authentication error handling

### Utils Tests (45.2% coverage)
- **utils_test.go**:
  - Tests password hashing and verification
  - Tests data hashing and consistency
  - Tests token generation with various lengths
  - Tests age validation with different dates
  - Tests password validation and matching

- **validation_test.go**:
  - Tests form validation for tenant applications
  - Tests conditional validation (e.g., if evicted, reason required)
  - Tests username and password generation
  - Tests various application form validations

## Setup Instructions

### Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Bevs-n-Devs/lilyshiddenparadise.git
   cd lilyshiddenparadise
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

   Required packages:
   - github.com/lib/pq - PostgreSQL driver
   - golang.org/x/crypto/bcrypt - Password hashing
   - github.com/joho/godotenv - Environment variable loading (optional)

### Configuration

1. Create a `.env` file in the `env` directory with the following variables:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_postgres_user
   DB_PASSWORD=your_postgres_password
   DB_NAME=lilyshiddenparadise
   MASTER_KEY=your_32_character_encryption_key
   ```

2. Set up the PostgreSQL database:
   ```bash
   psql -U postgres -c "CREATE DATABASE lilyshiddenparadise;"
   psql -U postgres -d lilyshiddenparadise -f db/schema.sql
   ```

### Running the Application

1. Build and run the application:
   ```bash
   go build -o lilyshiddenparadise
   ./lilyshiddenparadise
   ```

   Or use the Go run command:
   ```bash
   go run main.go
   ```

2. Access the application in your browser:
   ```
   http://localhost:8080
   ```

### Running Tests

1. Run all tests:
   ```bash
   go test ./...
   ```

2. Run tests with coverage:
   ```bash
   go test ./... -cover
   ```

3. Run tests for specific packages:
   ```bash
   go test ./middleware -v
   go test ./utils -v
   ```

4. Generate coverage report:
   ```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out
   ```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request
