# Build the CLI for the container's Linux architecture.
FROM golang:1.24-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/calculate-pi-cli .

# Runtime dependencies used by the calculation implementations.
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        default-jdk \
        python3 \
        python3-pip \
        tar \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://julialang-s3.julialang.org/bin/linux/x64/1.11/julia-1.11.2-linux-x86_64.tar.gz -o /tmp/julia.tar.gz \
    && tar -xzf /tmp/julia.tar.gz -C /usr/local --strip-components=1 \
    && rm /tmp/julia.tar.gz

WORKDIR /app
COPY --from=builder /out/calculate-pi-cli /app/calculate-pi-cli
RUN chmod 0755 /app/calculate-pi-cli

ENTRYPOINT ["/app/calculate-pi-cli"]
