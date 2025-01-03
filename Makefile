static:
	go vet -vettool=./cmd/staticlint ./...

run:
	go run ./cmd/shortener/main.go ./cmd/shortener

build:
	go build -o shortener ./cmd/shortener/*.go

test:
	go test ./...

fixmod:
	go mod tidy

proto:
	protoc --go_out=pkg/shortener_grpc  \
		--go-grpc_out=pkg/shortener_grpc \
		api/proto/shortener.proto

up:
	docker compose up -d

down:
	docker compose down