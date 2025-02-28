FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./internal/...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/main.go

FROM gcr.io/distroless/base-debian12

COPY --from=builder /app/main /main

CMD ["/main"]