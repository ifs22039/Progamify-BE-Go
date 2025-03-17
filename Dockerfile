# Menggunakan image Go resmi sebagai base image
FROM golang:1.20 AS builder

# Set working directory
WORKDIR /app

# Menyalin go.mod dan go.sum untuk mengelola dependensi
COPY go.mod go.sum ./
RUN go mod download

# Menyalin seluruh kode sumber ke dalam container
COPY . .

# Membangun aplikasi dari file main.go
RUN go build -o myapp ./cmd/api/main.go

# Menggunakan image minimal untuk menjalankan aplikasi
FROM alpine:latest

# Menyalin binary dari stage builder
COPY --from=builder /app/myapp .

# Menentukan port yang akan digunakan
EXPOSE 8080

# Menjalankan aplikasi
CMD ["./myapp"]