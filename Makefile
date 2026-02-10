# Database Connection String
DB_DSN := "postgres://postgres:password@localhost:5432/t_learn?sslmode=disable"

.PHONY: db-reset server seed stop

# 1. RESET DATABASE (Nuke -> Start -> Migrate)
db-reset: stop
	@echo "💥 Nuking database..."
	docker compose down -v
	@echo "🌱 Starting fresh database..."
	docker compose up -d db
	@echo "⏳ Waiting for database to be ready..."
	@sleep 3
	@echo "🏗️  Running Migrations..."
	goose -dir api/sql/schema postgres $(DB_DSN) up
	@echo "✅ Database is ready! Now run 'make server' in one terminal and 'make seed' in another."

# 2. RUN SERVER
server:
	@echo "🚀 Starting Server..."
	go run -C api cmd/server/main.go

# 3. RUN SEEDER
seed:
	@echo "🌳 Seeding Data..."
	go run -C api cmd/seeder/main.go

# Stop containers
stop:
	docker compose down
