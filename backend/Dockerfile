# syntax=docker/dockerfile:1

FROM golang:1.21-bookworm AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build

FROM golang:1.21-alpine3.19 AS runner

WORKDIR /app

COPY --from=builder /app/app ./

EXPOSE 8000

CMD ["/app/app"]

FROM golang:1.21-bookworm AS develop

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.49.0

CMD ["air"]
