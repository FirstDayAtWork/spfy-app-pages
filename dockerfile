# Define image to use
FROM golang:1.20.4

# ENV vars for hot reload
ENV PROJECT_DIR="/app" GO111MODULE="on" CGO_ENABLED=0
# Install a hot reload utility
RUN go install github.com/githubnemo/CompileDaemon@latest

# Set working dir of the container
WORKDIR /app

# Copy current directory contents to /app
COPY . /app

# Install the needed go deps
RUN go mod download

# Expose a port for the webapp to work
EXPOSE 2228

# Launch our app (first build then run)
ENTRYPOINT CompileDaemon -build="go build -o /mustracker_app" -command="/mustracker_app" -directory="/app" -exclude="**/*_test.go"
