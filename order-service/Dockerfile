FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
EXPOSE 3002
CMD ["go", "run", "main.go"]
