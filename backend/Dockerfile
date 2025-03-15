FROM golang:1.24

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Install SQLite
RUN apt-get update && apt-get install -y sqlite3

# Build the Go app
RUN CGO_ENABLED=1 go build -o ./main ./cmd/price-tracker

# Command to run the executable
CMD ["./main"]