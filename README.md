# t-learn Platform

t-learn is a full-stack educational platform designed to deliver interactive coding lessons. The system consists of a robust Go backend, a PostgreSQL database, a command-line interface (CLI) for student interaction, and a web frontend for progress tracking and administration.

## System Architecture

The platform is divided into three main components:

1. **Core API (Backend):** A RESTful API written in Go that manages authentication, content delivery (Courses, Lessons, Tasks), and student progress tracking.
2. **Database:** PostgreSQL database handling relational data with strict schema validation.
3. **Clients:**
   - **CLI:** A terminal-based client that fetches tasks and executes them in the user's local shell environment.
   - **Frontend:** A web dashboard for users to view course catalogs and for administrators to manage content.

## Directory Structure

```text
t-learn/
├── api/                # Backend API and CLI source code
│   ├── cmd/            # Entry points
│   │   ├── server/     # Main REST API server
│   │   └── seeder/     # Database population script
│   ├── internal/       # Private application logic (Auth, Content, DB)
│   └── sql/            # SQL queries and Goose migrations
└── frontend/           # Web application source code
```

## Technology Stack

- **Language:** Go (Golang) 1.22+
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **Database Migration:** Goose
- **SQL Generation:** SQLC
- **Authentication:** JWT (Web) and API Key (CLI)

## Getting Started

### Prerequisites

- Go 1.22 or higher

- Docker and Docker Compose

- Git

### Environment Configuration

Clone the repository and create the environment file.

```bash
git clone [https://github.com/yourusername/t-learn.git](https://github.com/yourusername/t-learn.git)
cd t-learn/api
cp .env.example .env
```

### Infrastructure Setup

Start the PostgreSQL container using Docker Compose.

```bash
docker compose up -d
```

### Database Migrations

Apply the database schema using Goose. This will create the users, courses, lessons, tasks, and completion tables.

```bash
# Install goose if not already installed
go install [github.com/pressly/goose/v3/cmd/goose@latest](https://github.com/pressly/goose/v3/cmd/goose@latest)

# Run migrations
cd sql/schema
goose postgres "postgres://postgres:password@localhost:5432/t_learn?sslmode=disable" up
cd ../..
```

### Seeding Data

To populate the database with initial courses and an admin user, run the seeder script. This script also serves as an integration test for the API.

```bash
go run cmd/seeder/main.go
```

### Running the Server

Start the API server.

```bash
go run cmd/server/main.go
```

The server will start on <http://localhost:8080>.

## API Documentation

### Authentication

- POST /auth/register - Register a new student account.

- POST /auth/login - Log in to receive a JWT.

- POST /auth/token - Generate a persistent API Key (used for CLI login).

### Public Content

- GET /courses - List all available courses.

- GET /courses/{id}/lessons - List all lessons for a specific course.

- GET /lessons/{id}/task - Fetch the task instructions and execution steps for a lesson.

### Student Actions

- POST /tasks/{id}/complete - Mark a task as completed (Requires Auth).

### Administration (Protected)

Requires a user with role='admin'.

- POST /admin/courses - Create a new course.

- POST /admin/courses/{id}/lessons - Add a lesson to a course.

- POST /admin/lessons/{id}/task - Create a multi-step task for a lesson.

- DELETE /admin/courses/{id} - Delete a course and all associated content.

## Frontend Setup

Instructions for setting up the frontend client.

1. Navigate to the frontend directory: cd frontend

2. Install dependencies: npm install

3. Start the development server: npm run dev

## CLI Integration

The t-learn CLI is maintained in a separate repository.

1. Install the CLI: Clone the CLI repository here: <https://github.com/Tikkaaa3/t-cli> and follow the build instructions.

2. Obtain API Key: Register an account via the Web Frontend or API, then request a key via POST /auth/token.

3. Authenticate:

```bash
t-cli login <YOUR_API_KEY>
```

## License

This project is licensed under the MIT License.

