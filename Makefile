name = sop112_exporter

all: build
.PHONY : all

test:
	go test ./...

lint:
	golint .

build: test lint
	mkdir -p builds
	go build -o ./builds/${name}
