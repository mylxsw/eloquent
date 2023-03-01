
run: generate-models
	go run _examples/main.go

generate-models: build
	./bin/eloquent gen --source "./_examples/models/*.yaml"
	go fmt ./_examples/models/*.go

build:
	go build -o bin/eloquent cmd/orm/*.go

init: build
	./bin/eloquent gen --source "./migrate/*.yaml"
	go fmt ./migrate/migrations.orm.go

.PHONY: build init generate-models run
