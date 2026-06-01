FROM ubuntu:22.04

# Install dependencies (fixed names)
RUN apt update && apt install -y \
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
    ca-certificates

# Build nsjail from source
RUN git clone https://github.com/google/nsjail.git /nsjail && \
    cd /nsjail && \
    make && \
    cp nsjail /usr/local/bin

# Install Go
RUN apt install -y golang

WORKDIR /app
COPY . .

RUN go build -o server ./cmd/goboxd

EXPOSE 8080
CMD ["./server"]