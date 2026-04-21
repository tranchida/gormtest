# syntax=docker/dockerfile:1

FROM docker.io/golang:1.26.2-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o gormtest .

FROM docker.io/alpine:3.21 AS runtime

ENV TZ=Europe/Zurich
RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=build /build/gormtest /app/gormtest
COPY --from=build /build/templates /app/templates

EXPOSE 8080

CMD ["/app/gormtest"]
