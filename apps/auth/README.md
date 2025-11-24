# Auth Service
 
 The Auth Service handles user registration, login, and JWT token generation/validation.
 
 ## 🏗️ Architecture
 
 ```mermaid
 sequenceDiagram
     participant User
     participant Frontend
     participant Auth
     participant DB
 
     User->>Frontend: Enter Credentials
     Frontend->>Auth: POST /login
     Auth->>DB: Query User
     DB-->>Auth: User Found
     Auth->>Auth: Verify Password (bcrypt)
     Auth->>Auth: Generate JWT
     Auth-->>Frontend: Return Token
     Frontend->>User: Login Success
 ```
 
 ## 🚀 API Endpoints
 
 | Method | Endpoint | Description |
 |--------|----------|-------------|
 | POST | `/register` | Register a new user |
 | POST | `/login` | Authenticate and get JWT |
 | GET | `/validate` | Validate JWT token |
 | GET | `/metrics` | Prometheus metrics |
 
 ## 🛠️ Tech Stack
 - **Framework**: Gin
 - **Database**: PostgreSQL
 - **ORM**: GORM
 - **Auth**: JWT (HS256)
