# Start from golang base image
FROM golang:alpine as builder

ENV POSTGRESQL_URL="postgres://gusampaio:gusampaio_pass@localhost:5432/hbday_db?sslmode=disable"

# Add Maintainer info
LABEL maintainer="Gustavo Sampaio"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git curl

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# migrating DB config
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-386.tar.gz | tar xvz
#RUN ./migrate -database ${POSTGRESQL_URL} -path db/migrations up

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]