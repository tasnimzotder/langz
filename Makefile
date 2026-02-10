.PHONY: build test vet install clean run docs

build:
	go build -o build/langz ./cmd/langz

test:
	gotestsum -- ./...

vet:
	go vet ./...

install:
	go install ./cmd/langz

clean:
	rm -rf build/

run:
	go run ./cmd/langz $(ARGS)

docs:
	mkdocs serve
