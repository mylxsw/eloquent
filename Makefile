
run: generate-models
	go run _examples/main.go

generate-models: build
	./bin/eloquent gen --source "./_examples/models/*.yml"
	go fmt ./_examples/models/*.go

build:
	go build -o bin/eloquent cmd/orm/*.go

init: build
	./bin/eloquent "./migrate/*.yml"
	go fmt ./migrate/migrations.orm.go

.PHONY: build init generate-models run
