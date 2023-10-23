FROM golang:1.20.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build ./cmd/kanastra

RUN chmod +x kanastra

EXPOSE 8080

CMD ["./kanastra"]