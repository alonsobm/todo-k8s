.PHONY: build

test:
	go test -v ./...

build:
	go build -o build/todo todoapp/cmd/main.go

run:
	go run todoapp/cmd/main.go

clean:
	rm -r ./build

pgdev:
	docker compose -f docker-compose.yml up -d database

mig.up:
	migrate -path=./migrations -database=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable up

mig.down:
	migrate -path=./migrations -database=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable down -all

test.integration:
	docker compose -f docker-compose.yml up -d database && sleep 5
	go test -v -run=TestIntegration ./todoapp/tests -count 1