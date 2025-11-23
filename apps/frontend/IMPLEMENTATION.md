# Frontend Service Implementation Plan

## Goal
Provide a user-friendly web interface for the PixelFlow application using React.

## Architecture
- **Type**: Single Page Application (SPA)
- **Framework**: React (Create React App)
- **Port**: 3000 (Dev Server)
- **Styling**: TailwindCSS
- **State Management**: React Context API

## Pages

### /login
- **Route**: Public
- **Features**: Email/Password form, JWT storage

### /register
- **Route**: Public
- **Features**: User registration form

### /dashboard
- **Route**: Protected (Requires JWT)
- **Features**:
  - Image URL upload form
  - List of user tasks
  - Real-time status updates (Polling)
  - Logout button

## Internal Components
- **AuthContext**: Manages user session and tokens
- **api**: Axios instance with interceptors for Authorization header
- **TaskList**: Component to render task cards
- **UploadForm**: Component to handle image submission

## Dependencies
- `react`, `react-dom`
- `react-router-dom` (Routing)
- `axios` (HTTP Client)
- `tailwindcss` (Styling)

## Configuration
- **Environment Variables**:
  - `REACT_APP_API_URL`: URL of the API Service (e.g., `http://localhost:8080`)
  - `REACT_APP_AUTH_URL`: URL of the Auth Service (e.g., `http://localhost:50051`)

## Deployment
- **Docker**: Multi-stage build
  1. **Build Stage**: Node.js image to compile React app (`npm run build`)
  2. **Production Stage**: Nginx Alpine image to serve static files
- **Nginx Config**: Custom `nginx.conf` to handle SPA routing (try_files $uri /index.html)
- **Orchestration**: Docker Compose service `frontend` mapped to port 3000
