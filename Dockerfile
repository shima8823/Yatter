# dev, builder
FROM golang:1.16 AS golang
WORKDIR /work/yatter-backend-go

# dev
FROM golang as dev
RUN mkdir -p /tmp/air-download && \
    curl -L -o /tmp/air-download/air.tar.gz https://github.com/cosmtrek/air/releases/download/v1.41.0/air_1.41.0_linux_amd64.tar.gz && \
    tar xf /tmp/air-download/air.tar.gz -C /tmp/air-download && \
    mv /tmp/air-download/air /go/bin/ && \
    chmod +x /go/bin/air && \
    rm -rf /tmp/air-download

# builder
FROM golang AS builder
COPY ./ ./
RUN make prepare build-linux

# release
FROM alpine AS app
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY --from=builder /work/yatter-backend-go/build/yatter-backend-go-linux-amd64 /usr/local/bin/yatter-backend-go
EXPOSE 8080
ENTRYPOINT ["yatter-backend-go"]
