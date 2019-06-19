
run: generate-models
	go run examples/main.go

generate-models: build
	./bin/orm "./examples/models/*.yml"
	go fmt ./examples/models/*.go

build:
	go build -o bin/orm cmd/orm/*.go