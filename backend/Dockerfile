FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod init
RUN go mod tidy

RUN go build -o main .

CMD ["./main"]