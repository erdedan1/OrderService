FROM golang:1.26

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 3001
CMD ["./app"]