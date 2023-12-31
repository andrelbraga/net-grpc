# syntax=docker/dockerfile:1
FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build cmd/main.go

EXPOSE 5001
CMD ["/build"]