# Подключаем переменные из .env (без кавычек и export!)
include .env
export

.PHONY: generate-server generate-client build clean print-env migrate-up migrate-down migrate-reset run stop

OPENAPI_GENERATOR=/home/linuxbrew/.linuxbrew/bin/openapi-generator
INPUT_FILE=open_api/swagger.yaml

generate-server:
	$(OPENAPI_GENERATOR) generate \
		-i $(INPUT_FILE) \
		-g go-gin-server \
		-o gen/server \
		--additional-properties=packageName=server,generateInterfaces=true \
		--skip-validate-spec --enable-post-process-file

generate-client:
	$(OPENAPI_GENERATOR) generate \
		-i $(INPUT_FILE) \
		-g go \
		-o gen/client \
		--additional-properties=packageName=client \
		--skip-validate-spec --enable-post-process-file

build: generate-server generate-client
	go build -o bin/pvz-service ./main/main.go

clean:
	rm -rf gen
	rm -rf bin

print-env:
	@echo "DSN_DB = $(DSN_DB)"
	@echo "TOKEN_MODERATOR = $(TOKEN_MODERATOR)"
	@echo "TOKEN_CLIENT = $(TOKEN_CLIENT)"

# Команды для управления миграциями
migrate-up:
	@echo "🚀 Применение миграций..."
	go run cmd/migrate/main.go -up

migrate-down:
	@echo "⚠️  Удаление всех таблиц и данных..."
	go run cmd/migrate/main.go -down

migrate-reset: migrate-down migrate-up
	@echo "🔄 База данных сброшена и пересоздана"

# Команды для запуска и остановки приложения
run: migrate-up
	@echo "🚀 Запуск приложения..."
	go run main/main.go

stop:
	@echo "🛑 Остановка приложения..."
	@pkill -f "go run main/main.go" || true
	@pkill -f "pvz-service" || true
