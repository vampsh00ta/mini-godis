FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ../../Downloads/avito_intern_spring-main%20 ./
RUN CGO_ENABLED=0 GOOS=linux go build  ./cmd/mini-godis.
CMD ["./avito_intern"]
