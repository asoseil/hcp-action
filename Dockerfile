# Compact image
FROM golang:1.23.3-alpine3.20

# Set wd
WORKDIR /app

# Copy repo content
COPY ./ ./

# Install module and deps
RUN go mod download
RUN go mod tidy

# TODO move away the build to speed up the startup
# Create bin file
RUN go build -o /bin/app .

# Showtime
ENTRYPOINT ["/bin/app"]