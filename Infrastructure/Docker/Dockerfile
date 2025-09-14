FROM golang:1.21 AS build

WORKDIR /app

COPY .netrc /root/.netrc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mcabank-auth cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/mcabank-auth .

EXPOSE 8080

CMD ["./mcabank-auth"]

