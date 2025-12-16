FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем main.go из cmd/app
RUN go build -o wallet_app ./cmd/app

EXPOSE 8080

CMD ["./wallet_app"]
