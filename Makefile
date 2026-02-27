.PHONY: build test deploy dev dev-up dev-down clean

build:
	mkdir -p deploy
	GOOS=linux GOARCH=arm64 go build -o deploy/bootstrap cmd/lambda/main.go

test:
	go test ./...

deploy: build
	cd infra && cdk deploy --require-approval never

dev-up:
	docker compose up -d

dev-down:
	docker compose down

dev: dev-up
	go run cmd/local/main.go

clean:
	rm -rf deploy
