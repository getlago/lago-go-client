FROM golang:1.24

RUN apt-get update && apt-get install -y build-essential git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && \
  go mod tidy && \
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.0

COPY . .

CMD ["go", "test", "./..."]
