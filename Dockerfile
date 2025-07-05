# Stage 1: Build the application
FROM golang:latest as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o wom_backend

# Stage 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/wom_backend .

# Create directories that will be volume mounted
RUN mkdir -p /root/elections /root/svgs

# Command to run the executable
CMD ["./wom_backend"]