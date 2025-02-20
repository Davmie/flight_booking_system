FROM golang:1.22

WORKDIR /

COPY . /app

WORKDIR /app

RUN go mod tidy
RUN go mod download
RUN go build cmd/main.go

EXPOSE 8050

CMD ["./main"]
