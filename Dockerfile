# use : docker build -t cosmo-proof:1.0.0 .
FROM golang:1.14.2 as builder

RUN mkdir -p /go/src/github.com/Alexader/
WORKDIR /go/src/github.com/meshplus/pier

# Cache dependencies
COPY go.mod .
COPY go.sum .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download -x

# Build real binaries
COPY . .

RUN make install

# Final image
FROM frolvlad/alpine-glibc

WORKDIR /root

# Copy over binaries from the builder
COPY --from=builder /go/bin/proof /usr/local/bin

RUN ["proof", "start"]

EXPOSE 44555 44544