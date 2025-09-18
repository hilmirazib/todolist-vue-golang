FROM golang:1.24.0

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .


CMD ["go", "run", "main.go"]