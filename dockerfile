# Define image to use
FROM golang:1.20.4

# Set working dir of the container
WORKDIR /app

# Copy current directory contents to /app
COPY . /app

# Install the needed go deps
RUN go mod download

# Build go app
RUN go build -o /mustracker_app

# Expose a port for the webapp to work
EXPOSE 2228

# Define an ENV variable JIC
ENV NAME Ciao

# Launch our app
CMD [ "/mustracker_app" ]
