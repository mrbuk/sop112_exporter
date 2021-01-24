FROM golang:1.15-alpine AS build

WORKDIR /go/src/sop112_exporter
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -v ./...

# build a minimal image

FROM scratch

# the tls certificates:
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# the actual binary
COPY --from=build /go/bin/sop112_exporter /go/bin/sop112_exporter

ENTRYPOINT ["/go/bin/sop112_exporter"] 
