MAKEFLAGS += --no-print-directory

BINARY_NAME=myapp
# DSN="host=localhost port=5432 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5"
DSN="user=postgres password=password host=localhost port=5432 dbname=concurrency sslmode=disable"
REDIS="127.0.0.1:6379"

start-compose:
	@echo "Starting compose..."
	@docker compose up -d
	@echo "Started!"

stop-compose:
	@echo "Stopping compose..."
	@docker compose down
	@echo "Stoped!"

build:
	@echo "building..."
	@env CGO_ENABLED=0 go build -ldflags="-s -w" -o ../go-out/${BINARY_NAME} ../cmd/web
	@echo "built!"

run: start-compose build
	@echo "Starting"
	@env DSN=${DSN} REDIS=${REDIS} ../go-out/${BINARY_NAME} &
	@echo "Started"

clean:
	@echo "Cleaning..."
	@go clean
	@rm ../go-out/${BINARY_NAME}
	@echo "Cleaned!"

start: run

stop: stop-compose
	@echo "Stopping..."
	# @-pkill -SIGTERM -f "${BINARY_NAME}" 2>/dev/null || true
	@-pkill "${BINARY_NAME}" 2>/dev/null || true
	@echo "Stoped!"

restart: stop start

test:
	go test -v ../...
