FROM golang:1.22.5-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod vendor

COPY . .

RUN go build -o ./main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
