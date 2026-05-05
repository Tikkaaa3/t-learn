# t-learn Platform

<p align="center">

![License](https://img.shields.io/github/license/Tikkaaa3/t-learn?style=for-the-badge)
![Last Commit](https://img.shields.io/github/last-commit/Tikkaaa3/t-learn?style=for-the-badge)
![Repo Size](https://img.shields.io/github/repo-size/Tikkaaa3/t-learn?style=for-the-badge)

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-blue?style=for-the-badge&logo=typescript)
![React](https://img.shields.io/badge/React-19.2-61DAFB?style=for-the-badge&logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=for-the-badge&logo=docker)

</p>

**t-learn** is a modern, full-stack educational platform designed to deliver interactive coding lessons through a terminal-based interface. The system combines a robust Go backend, PostgreSQL database, and a React-based web terminal frontend to provide an immersive learning experience.

## 🎬 Demo Video

👉 Click the image below to watch the demo on YouTube:

[![▶ Watch Demo](https://img.youtube.com/vi/IJQU1nXh-sI/maxresdefault.jpg)](https://www.youtube.com/watch?v=IJQU1nXh-sI)

## 🌟 Features

### For Students

- **Terminal-Based Learning**: Learn through an interactive web terminal with a Unix-like command interface
- **Progress Tracking**: Automatic tracking of completed lessons and tasks
- **Course Navigation**: Browse courses and lessons using familiar `cd` and `ls` commands
- **Multi-Step Tasks**: Follow structured, step-by-step coding exercises
- **Persistent Sessions**: Auto-login functionality for seamless experience across sessions
- **Rich Markdown Content**: Beautiful rendering of lesson content with code blocks and formatting

### For Administrators

- **Course Management**: Create, update, and delete courses through admin commands
- **Lesson Creation**: Build structured lessons with multiple steps and verification tasks
- **Role-Based Access**: Secure admin endpoints with role-based authentication
- **Content Seeding**: Quick database population for testing and demos

### Terminal Features

- **Full Cursor Support**: Navigate text with arrow keys, Home, End, Delete
- **Command History**: Use ↑/↓ to browse previous commands
- **Auto-completion Ready**: Extensible command system
- **Copy/Paste Support**: Standard clipboard operations work seamlessly

## 🏗️ System Architecture

The platform consists of three main components:

1. **Backend API** (Go): RESTful API handling authentication, content delivery, and progress tracking
2. **Database** (PostgreSQL): Relational database with strict schema validation
3. **Web Frontend** (React + TypeScript): Interactive terminal UI with directory-based navigation

```
┌─────────────────┐
│  Web Frontend   │  React + TypeScript + Vite
│  (Terminal UI)  │  Terminal emulation with Markdown rendering
└────────┬────────┘
         │ HTTPS/REST
         ▼
┌─────────────────┐
│   Backend API   │  Go + net/http
│  (Port 8080)    │  JWT Auth + API Keys
└────────┬────────┘
         │ SQL
         ▼
┌─────────────────┐
│   PostgreSQL    │  Persistent storage
│  (Port 5432)    │  SQLC + Goose migrations
└─────────────────┘
```

## 📁 Directory Structure

```
t-learn/
├── api/                      # Backend (Go)
│   ├── cmd/
│   │   ├── server/          # Main API server entry point
│   │   └── seeder/          # Database seeding script
│   ├── internal/
│   │   ├── auth/            # Authentication (JWT, API keys, password hashing)
│   │   ├── content/         # Course/Lesson/Task handlers
│   │   └── database/        # SQLC generated code + models
│   ├── sql/
│   │   ├── queries/         # SQLC query definitions
│   │   └── schema/          # Goose migrations
│   ├── go.mod
│   └── .env.example
├── frontend/                 # Frontend (React + TypeScript)
│   ├── src/
│   │   ├── api/             # API client functions
│   │   ├── commands/        # Terminal command registry
│   │   ├── hooks/           # React hooks (useTerminal)
│   │   ├── App.tsx          # Main terminal component
│   │   └── types.ts         # TypeScript interfaces
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml        # PostgreSQL container
├── Makefile                  # Build automation
└── README.md
```

## 🛠️ Technology Stack

### Backend

- **Language**: Go 1.25+
- **Web Framework**: Native `net/http` with custom routing
- **Database**: PostgreSQL 15
- **SQL Toolkit**:
  - [SQLC](https://sqlc.dev/) - Type-safe SQL code generation
  - [Goose](https://github.com/pressly/goose) - Database migrations
- **Authentication**:
  - JWT (golang-jwt/jwt) for web sessions
  - API Keys for CLI access
  - Argon2id for password hashing
- **Dependencies**:
  - `jackc/pgx/v5` - PostgreSQL driver
  - `google/uuid` - UUID generation
  - `joho/godotenv` - Environment configuration

### Frontend

- **Framework**: React 19.2
- **Language**: TypeScript 5.9
- **Build Tool**: Vite 7.3
- **Styling**: CSS with CSS variables
- **Markdown**: react-markdown for rich content rendering
- **State Management**: React hooks (useState, useCallback, useEffect)

### Infrastructure

- **Containerization**: Docker Compose
- **Database**: PostgreSQL 15 (Docker image)
- **Development**: Hot reload with Vite dev server

## 🚀 Getting Started

### Prerequisites

- **Go** 1.25 or higher ([Download](https://go.dev/dl/))
- **Node.js** 18+ and npm ([Download](https://nodejs.org/))
- **Docker** and Docker Compose ([Download](https://www.docker.com/))
- **Git**

### Quick Start (Using Makefile)

The easiest way to get started:

```bash
# 1. Clone the repository
git clone https://github.com/Tikkaaa3/t-learn.git
cd t-learn

# 2. Reset database (starts PostgreSQL, runs migrations)
make db-reset

# 3. Start the backend server (in one terminal)
make server

# 4. Seed the database with sample data (in another terminal)
make seed

# 5. Start the frontend (in a third terminal)
cd frontend
npm install
npm run dev
```

Visit **<http://localhost:5173>** to access the terminal interface.

### Manual Setup (Step by Step)

#### 1. Environment Configuration

```bash
cd api
cp .env.example .env
```

Edit `.env` with your configuration:

```env
PORT=8080
DB_URL=postgres://postgres:password@localhost:5432/t_learn?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-here
```

#### 2. Start PostgreSQL

```bash
docker compose up -d
```

This starts PostgreSQL on port 5432.

#### 3. Install Goose (Migration Tool)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

#### 4. Run Database Migrations

```bash
cd api/sql/schema
goose postgres "postgres://postgres:password@localhost:5432/t_learn?sslmode=disable" up
cd ../../..
```

This creates 6 tables:

- `users` - User accounts and authentication
- `courses` - Course catalog
- `lessons` - Lesson content within courses
- `tasks` - Multi-step tasks for each lesson
- `steps` - Individual steps within tasks
- `task_completions` - User progress tracking

#### 5. Seed the Database

```bash
cd api
go run cmd/seeder/main.go
```

This creates:

- An admin user: `admin` / `adminpass`
- Sample courses (Go Mastery, Rust Basics, etc.)
- Lessons with multi-step tasks

#### 6. Start the Backend Server

```bash
cd api
go run cmd/server/main.go
```

Server runs on **<http://localhost:8080>**

#### 7. Start the Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on **<http://localhost:5173>**

## 🖥️ Using the Terminal Interface

### First Time Setup

1. **Register an account**:

   ```bash
   $ register myusername my@email.com mypassword
   ```

2. **Login**:

   ```bash
   $ login myusername mypassword
   ```

3. **Check your profile**:

   ```bash
   $ whoami
   # Shows: username and completed tasks count
   ```

### Navigating Courses

The terminal uses a directory-based navigation system similar to Unix:

```bash
# See available commands
$ help

# Navigate to courses directory
$ cd courses
courses $

# List all available courses
courses $ ls
📁 Go Mastery/
📁 Rust Basics/
📁 Python Fundamentals/

# Enter a specific course
courses $ cd "Go Mastery"
courses/Go Mastery $

# List lessons in the course
courses/Go Mastery $ ls
[ ] 📄 Introduction to Go
[✓] 📄 Variables and Types
[ ] 📄 Functions and Methods

# Start a lesson
courses/Go Mastery $ start "Introduction to Go"

# Go back to courses directory
courses/Go Mastery $ cd ..
courses $

# Go back to root
courses $ cd ..
$
```

### Terminal Commands Reference

| Command                          | Description                 | Example                                  |
| -------------------------------- | --------------------------- | ---------------------------------------- |
| `help`                           | Show all available commands | `help`                                   |
| `clear`                          | Clear the terminal screen   | `clear`                                  |
| `register <user> <email> <pass>` | Create a new account        | `register john john@example.com pass123` |
| `login <user> <pass>`            | Login to your account       | `login john pass123`                     |
| `logout`                         | Logout and clear session    | `logout`                                 |
| `whoami`                         | Show username and stats     | `whoami`                                 |
| `token`                          | Generate API key for CLI    | `token`                                  |
| `cd <dir>`                       | Change directory            | `cd courses`                             |
| `cd ..`                          | Go back one level           | `cd ..`                                  |
| `ls`                             | List courses or lessons     | `ls`                                     |
| `start <lesson>`                 | Start a lesson task         | `start "Intro to Go"`                    |

### Admin Commands

For users with `role='admin'`:

| Command                                         | Description     | Example                                          |
| ----------------------------------------------- | --------------- | ------------------------------------------------ |
| `mkcourse "<title>" "<desc>"`                   | Create a course | `mkcourse "Go Mastery" "Learn Go"`               |
| `rmcourse <course>`                             | Delete a course | `rmcourse "Go Mastery"`                          |
| `mklesson <course> <pos> "<title>" "<content>"` | Create a lesson | `mklesson "Go Mastery" 1 "Intro" "Content here"` |
| `rmlesson <lesson>`                             | Delete a lesson | `rmlesson "Intro"`                               |

### Keyboard Shortcuts

- **Arrow Up/Down** (↑/↓): Navigate command history
- **Arrow Left/Right** (←/→): Move cursor within input
- **Home**: Jump to start of input
- **End**: Jump to end of input
- **Backspace**: Delete character before cursor
- **Delete**: Delete character after cursor
- **Enter**: Execute command
- **Ctrl+V**: Paste text at cursor position

## 📡 API Documentation

### Base URL

`http://localhost:8080`

### Authentication Endpoints

#### Register

```http
POST /auth/register
Content-Type: application/json

{
  "username": "john",
  "email": "john@example.com",
  "password": "securepass"
}
```

#### Login

```http
POST /auth/login
Content-Type: application/json

{
  "username": "john",
  "password": "securepass"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid-here",
    "username": "john"
  }
}
```

#### Generate API Key (Protected)

```http
POST /auth/token
Authorization: Bearer <jwt_token>

Response:
{
  "api_key": "generated-api-key-here"
}
```

#### Get User Stats (Protected)

```http
GET /auth/stats
Authorization: Bearer <jwt_token>

Response:
{
  "username": "john",
  "completed_tasks": 5
}
```

### Content Endpoints

#### List Courses

```http
GET /courses

Response:
[
  {
    "id": "uuid",
    "title": "Go Mastery",
    "description": "Learn Go programming",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### List Lessons (Protected)

```http
GET /courses/{course_id}/lessons
Authorization: Bearer <jwt_token>

Response:
[
  {
    "id": "uuid",
    "course_id": "uuid",
    "title": "Introduction to Go",
    "content": "# Lesson content in markdown",
    "position": 1,
    "completed": false
  }
]
```

#### Get Task

```http
GET /lessons/{lesson_id}/task

Response:
{
  "lesson_id": "uuid",
  "lesson_title": "Introduction to Go",
  "lesson_content": "# Content",
  "task_id": "uuid",
  "task_description": "Create a Hello World program",
  "steps": [
    {
      "position": 1,
      "command": "mkdir hello-world && cd hello-world"
    }
  ]
}
```

#### Complete Task (Protected)

```http
POST /tasks/{task_id}/complete
Authorization: Bearer <jwt_token>

Response:
{
  "message": "Task marked as complete!"
}
```

### Admin Endpoints

All admin endpoints require `role='admin'` in the user record.

#### Create Course

```http
POST /admin/courses
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Go Mastery",
  "description": "Master Go programming"
}
```

#### Create Lesson

```http
POST /admin/courses/{course_id}/lessons
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Introduction",
  "content": "# Markdown content here",
  "position": 1
}
```

#### Create Task

```http
POST /admin/lessons/{lesson_id}/task
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "description": "Complete the exercise",
  "steps": [
    {"position": 1, "command": "mkdir project"},
    {"position": 2, "command": "cd project"}
  ]
}
```

#### Delete Course

```http
DELETE /admin/courses/{course_id}
Authorization: Bearer <jwt_token>
```

#### Delete Lesson

```http
DELETE /admin/lessons/{lesson_id}
Authorization: Bearer <jwt_token>
```

#### Delete Task

```http
DELETE /admin/tasks/{task_id}
Authorization: Bearer <jwt_token>
```

## 🔧 Development

### Backend Development

```bash
cd api

# Install dependencies
go mod download

# Generate SQL code (after modifying queries)
sqlc generate

# Run tests
go test ./...

# Build
go build -o server cmd/server/main.go

# Run with hot reload (install air first)
air
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Development server with hot reload
npm run dev

# Type checking
npm run build

# Linting
npm run lint

# Production build
npm run build
npm run preview
```

### Database Management

```bash
# Create a new migration
cd api/sql/schema
goose create add_new_table sql

# Run migrations
goose postgres $DB_URL up

# Rollback last migration
goose postgres $DB_URL down

# Check migration status
goose postgres $DB_URL status
```

### Environment Variables

**Backend** (`.env` in `api/` directory):

```env
PORT=8080
DB_URL=postgres://postgres:password@localhost:5432/t_learn?sslmode=disable
JWT_SECRET=your-super-secret-key-change-in-production
```

**Frontend** (Vite automatically picks up `VITE_` prefixed vars):
The API URL is currently hardcoded in `frontend/src/api/client.ts` as `http://localhost:8080`

## 🔐 Security Features

- **Password Hashing**: Argon2id algorithm for secure password storage
- **JWT Authentication**: Stateless authentication with configurable secrets
- **API Keys**: Persistent tokens for CLI/external integrations
- **Role-Based Access Control**: Admin-only endpoints protected by middleware
- **CORS Protection**: Configured to allow frontend origin only
- **SQL Injection Protection**: Parameterized queries via SQLC
- **Auto-Login**: Secure session persistence using localStorage

## 🐳 Docker Deployment

### Using Docker Compose (Recommended)

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Reset (remove volumes)
docker compose down -v
```

### Manual Docker Commands

```bash
# PostgreSQL only
docker run -d \
  --name t-learn-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=t_learn \
  -p 5432:5432 \
  postgres:15
```

## 📊 Database Schema

The platform uses 6 core tables:

### `users`

- User accounts with authentication credentials
- Fields: `id`, `username`, `email`, `password_hash`, `role`, `api_key`

### `courses`

- Course catalog
- Fields: `id`, `title`, `description`, `created_at`, `updated_at`

### `lessons`

- Lessons within courses
- Fields: `id`, `course_id`, `title`, `content`, `position`

### `tasks`

- Tasks/exercises for lessons
- Fields: `id`, `lesson_id`, `description`, `created_at`, `updated_at`

### `steps`

- Individual steps within tasks
- Fields: `id`, `task_id`, `position`, `command`

### `task_completions`

- Tracks which users completed which tasks
- Fields: `id`, `user_id`, `task_id`, `created_at`, `updated_at`

## 🧪 Testing

### Seeder as Integration Test

The seeder (`cmd/seeder/main.go`) serves as both a data population tool and an integration test:

```bash
go run cmd/seeder/main.go
```

It tests:

- User registration and login
- Admin role assignment
- Course creation
- Lesson creation with tasks
- Task completion workflow

### Manual Testing

1. **Registration Flow**:

   ```bash
   $ register testuser test@example.com testpass
   $ login testuser testpass
   $ whoami
   ```

2. **Navigation Flow**:

   ```bash
   $ cd courses
   courses $ ls
   courses $ cd "Go Mastery"
   courses/Go Mastery $ ls
   ```

3. **Admin Flow** (requires admin role):

   ```bash
   $ mkcourse "New Course" "Description"
   $ cd courses
   courses $ cd "New Course"
   courses/New Course $ ls
   ```

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- **Go**: Follow standard Go formatting (`gofmt`, `golint`)
- **TypeScript**: Use ESLint configuration provided
- **Git**: Write clear, descriptive commit messages

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 Related Projects

- **t-cli**: Command-line client for t-learn platform ([GitHub](https://github.com/Tikkaaa3/t-cli))

## 📧 Support

For issues, questions, or suggestions:

- Open an issue on GitHub
- Contact: <tikkaaa3@gmail.com>

## 🗺️ Roadmap

Future planned features:

- [ ] Progress analytics dashboard
- [ ] Social features (leaderboards, sharing)
- [ ] Mobile app
- [ ] Content marketplace
