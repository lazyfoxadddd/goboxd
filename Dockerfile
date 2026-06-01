FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    pkg-config \
    libprotobuf-dev \
    protobuf-compiler \
    libnl-3-dev \
    libnl-route-3-dev \
    python3 \
    gcc \
    curl \
    ca-certificates \
    golang \
    && rm -rf /var/lib/apt/lists/*

# Build nsjail
RUN git clone https://github.com/google/nsjail.git /nsjail && \
    cd /nsjail && \
    make && \
    cp nsjail /usr/local/bin && \
    chmod +x /usr/local/bin/nsjail

WORKDIR /app
COPY . .

RUN go build -o goboxd ./cmd/goboxd

EXPOSE 8080
CMD ["./goboxd"]