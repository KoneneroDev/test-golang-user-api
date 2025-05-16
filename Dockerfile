FROM golang:1.24.2

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/test-golang-user-api
RUN go build -o /app/app .
WORKDIR /app
EXPOSE 8080
CMD ["./app"]

