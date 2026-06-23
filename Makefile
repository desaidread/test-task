include .env
export

migrate-up:
	migrate -path ./migrations -database "${SQL_CONN}" up

migrate-down:
	migrate -path ./migrations -database "${SQL_CONN}" down