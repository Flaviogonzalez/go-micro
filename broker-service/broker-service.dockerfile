FROM golang:1.24-alpine AS builder

# creo el container /app
RUN mkdir /app

# copy all the files from the current directory to /app in the container
COPY . /app

# set the Current Working Directory inside the container
WORKDIR /app

# install dependencies
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# creates a self-contained binary
RUN chmod +x /app/brokerApp

# build a tiny broker image
FROM alpine:latest

RUN mkdir /app 

# copy the binary from the builder stage to the new image
COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]