.PHONY: build test deploy dev clean

build:
	mkdir -p deploy
	GOOS=linux GOARCH=arm64 go build -o deploy/bootstrap cmd/lambda/main.go

test:
	go test ./...

deploy: build
	cd infra && cdk deploy --require-approval never

dev:
	go run cmd/local/main.go

clean:
	rm -rf deploy
