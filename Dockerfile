FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instala dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o load-tester cmd/main.go

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/load-tester .

ENTRYPOINT ["./load-tester"]
