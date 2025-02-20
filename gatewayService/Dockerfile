FROM golang:1.22

WORKDIR /

COPY . /app


WORKDIR /app/gatewayService

RUN go mod tidy
RUN go mod download
RUN go build cmd/main.go

EXPOSE 8080

CMD ["./main"]
