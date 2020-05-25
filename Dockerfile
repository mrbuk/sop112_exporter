FROM golang:1.14

WORKDIR /go/src/sop112_exporter
COPY . .

RUN go mod tidy
RUN go install -v ./...

CMD ["sh", "-c", "sop112_exporter -broadcast=${BCAST_ADDRESS}"]
