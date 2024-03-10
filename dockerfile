# Define image to use
FROM golang:1.20.4

# Installing a text editor JIC
ENV USE_POLLING_FILE_WATCHER=true
RUN apt-get update && apt-get install nano

RUN go install github.com/cortesi/modd/cmd/modd@latest

# Set working dir of the container
WORKDIR /app

# Copy current directory contents to /app
COPY . /app

# Install the needed go deps
RUN go mod download

# Expose a port for the webapp to work
EXPOSE 2228

# Launch our app (first build then run)
CMD ["modd"]