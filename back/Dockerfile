FROM golang:1.20.5 as build
WORKDIR /opt/sandbox-docker-go
COPY . .
RUN go build -o app main.go
CMD ["/opt/sandbox-docker-go/app"]