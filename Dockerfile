FROM golang:1.25 AS build

WORKDIR /build

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -o app

FROM ubuntu:24.04

WORKDIR /app

COPY --from=build /build/app .
COPY --from=build /build/repo/migrations ./migrations

CMD ["./app"]
