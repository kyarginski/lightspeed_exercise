lint:
	golangci-lint run ./...

build_counter:
	go build -o bin/counter cmd/counter/main.go

build_generator:
	go build -o bin/generator cmd/generator/main.go

run:
	bin/counter result.txt 1