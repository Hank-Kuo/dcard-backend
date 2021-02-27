# Golang Server's Dockerfile
FROM golang:1.14-alpine

ENV GIN_MODE=release

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Add Maintainer Info
LABEL maintainer="Hank Kuo <asdf024681029@gmail.com>"

WORKDIR /docker-backend

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o server ./cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
ENTRYPOINT ["./server"]