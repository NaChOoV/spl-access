ent-generate:
	ent generate ./ent/schema --feature sql/upsert --feature sql/execquery

build:
	go build -o main main.go

run:
	go run main.go

test:
	go test -v ./src/service/... ./src/controller/... ./src/helpers/...