FROM golang:alpine AS builder
ENV CGO_ENABLED 0
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -ldflags="-s -w" -o load28 .
FROM scratch
WORKDIR /build
COPY --from=builder /build/load28 /build/load28
CMD ["./load28"]