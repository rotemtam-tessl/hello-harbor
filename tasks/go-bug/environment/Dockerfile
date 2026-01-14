FROM golang:1.24

WORKDIR /workspace

# Install SQLite, build tools for CGO, and Atlas CLI
RUN apt-get update && apt-get install -y gcc sqlite3 curl && rm -rf /var/lib/apt/lists/*
RUN curl -sSf https://atlasgo.sh | sh

# Create logs directory for verifier output
RUN mkdir -p /logs/verifier

# Copy Go module and download dependencies
COPY go.mod ./
RUN go mod tidy && go mod download
