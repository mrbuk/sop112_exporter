name = sop112_exporter

all: build
.PHONY : all

run-integrationtest:
	ginkgo -r

test:
	ginkgo -r -skipPackage integrationtest
	
lint:
	golint .
	
build: test lint
	mkdir -p build
	go build -o ./build/${name}
