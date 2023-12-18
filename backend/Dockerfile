# syntax=docker/dockerfile:1

FROM golang:1.21.5-bookworm AS builder

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build

FROM golang:1.21.5-alpine AS runner

WORKDIR /app

COPY --from=builder /app/app ./

EXPOSE 8000

CMD ["/app/app"]

FROM golang:1.21.5 AS develop

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.49.0

CMD ["air"]
