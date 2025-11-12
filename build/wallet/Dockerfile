FROM golang:1.24-alpine AS build

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/wallet ./cmd
COPY internal ./internal

RUN GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd

FROM alpine:latest AS app

WORKDIR /app

COPY --from=build /app/bin/main .

CMD ["./main"]
