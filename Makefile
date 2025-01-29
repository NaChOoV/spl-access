ent-generate:
	ent generate ./ent/schema --feature sql/upsert --feature sql/execquery

build:
	go build -o bin/myapp main.go

run:
	go run main.go