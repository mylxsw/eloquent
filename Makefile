
run: generate-models
	go run _examples/main.go

generate-models: build
	./bin/orm "./_examples/models/*.yml"
	go fmt ./_examples/models/*.go

build:
	go build -o bin/orm cmd/orm/*.go

init: build
	./bin/orm "./migrate/*.yml"
	go fmt ./migrate/migrations.orm.go

.PHONY: build init generate-models run
