FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
COPY ./database /app/databaseы
RUN go build -o main ./cmd

CMD ["./main"]