FROM golang:1.20.5 as build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build ./main.go

EXPOSE 8080

CMD ["./main"]
