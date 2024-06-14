FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o mdns-tcp-proxy .

FROM scratch

WORKDIR /app

COPY --from=builder /app/mdns-tcp-proxy ./mdns-tcp-proxy

ENTRYPOINT ["./mdns-tcp-proxy"]
