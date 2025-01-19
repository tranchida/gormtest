# syntax=docker/dockerfile:1

FROM docker.io/golang:1.23.4-alpine AS build

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o gormtest cmd/gormtest/main.go

FROM docker.io/alpine:3.14 AS runtime

ENV TZ=Europe/Zurich
RUN apk --no-cache add tzdata

COPY --from=build /build/gormtest /app/gormtest

EXPOSE 8080

CMD ["/app/gormtest"]
