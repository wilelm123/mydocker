FROM golang:1.21-alpine AS builder

LABEL maintainer="leich3(leich3@cisco.com)"

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm64
RUN go build -ldflags="-s -w" -o mydocker .

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/mydocker", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/mydocker"]
