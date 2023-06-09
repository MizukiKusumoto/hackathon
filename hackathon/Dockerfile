FROM golang:1.20.5 as build
WORKDIR /opt/sandbox-docker-go
COPY . .
RUN go mod download github.com/go-sql-driver/mysql
RUN go mod download github.com/oklog/ulid
RUN go build -o app main.go
CMD ["/opt/sandbox-docker-go/app"]