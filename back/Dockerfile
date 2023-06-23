FROM golang:1.20.5 as build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build ./main.go

CMD ["./main"]