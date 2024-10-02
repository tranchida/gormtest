# syntax=docker/dockerfile:1

FROM debian:stable-slim

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY bin/gormtest .
COPY sample.jpg .

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

ENV GIN_MODE=release
# Run
CMD ["./gormtest"]
