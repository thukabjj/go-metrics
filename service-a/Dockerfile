FROM golang:1.19 as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY .. .
RUN go build -o metrics

ENTRYPOINT ["/app/metrics"]