FROM golang:1.17.3-alpine3.14

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o latency .

CMD ["./latency"]