name = sop112_exporter

all: build
.PHONY : all

test:
	ginkgo -r

lint:
	golint .

build: test lint
	mkdir -p build
	go build -o ./build/${name}
