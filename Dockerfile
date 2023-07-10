# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /build

# Copy go module dependencies into the container
COPY go.mod go.sum .

# Download dependencies
RUN go mod download

# Copy go code into the container
COPY . .

# Copy the ports.json file to the build context
COPY ports.json .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

# Production stage
FROM scratch

# Copy the binary from the build stage to the production image
COPY --from=builder /build/app ./app/server

# Copy the ports.json file to the production image
COPY ports.json /app/ports.json

# Use an unprivileged user.
USER ${USER}:${USER}

WORKDIR /app

# Set the entry point to run the application
ENTRYPOINT ["./server"]
