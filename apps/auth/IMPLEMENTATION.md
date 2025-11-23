# Auth Service Implementation Plan

## Goal
Provide secure user authentication and management using HTTP REST APIs.

## Architecture
- **Type**: REST API
- **Port**: 50051
- **Database**: PostgreSQL (User data)
- **Auth Mechanism**: JWT (JSON Web Tokens)
- **Password Hashing**: Bcrypt

## API Endpoints

### POST /register
- **Input**: `email`, `password`
- **Logic**:
  1. Validate input
  2. Hash password with bcrypt
  3. Create user in PostgreSQL
  4. Return success

### POST /login
- **Input**: `email`, `password`
- **Logic**:
  1. Find user by email
  2. Compare password hash
  3. Generate JWT token (valid for 24h)
  4. Return token

### GET /validate
- **Input**: `Authorization: Bearer <token>`
- **Logic**:
  1. Parse JWT token
  2. Verify signature and expiration
  3. Return user ID

## Internal Components
- **db**: PostgreSQL connection using GORM
- **utils**: JWT generation and validation, Password hashing
- **models**: User struct definition

## Dependencies
- `github.com/gin-gonic/gin` (HTTP Framework)
- `gorm.io/gorm` (ORM)
- `gorm.io/driver/postgres` (DB Driver)
- `github.com/golang-jwt/jwt/v5` (Tokens)
- `golang.org/x/crypto/bcrypt` (Hashing)
