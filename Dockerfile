# Menggunakan image Go resmi sebagai base image
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Menyalin go.mod dan go.sum untuk mengelola dependensi
COPY .env /app/.env
COPY go.mod go.sum ./
RUN go mod download

# Menyalin seluruh kode sumber ke dalam container
COPY . .

# Membangun aplikasi secara statis
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp ./cmd/api/main.go

# Menggunakan image minimal untuk menjalankan aplikasi
FROM ubuntu:latest

# Set working directory
WORKDIR /app

# Menyalin binary dari stage builder
COPY --from=builder /app/myapp /app/myapp
COPY --from=builder /app/.env /app/.env

# Memberikan izin eksekusi
RUN chmod +x /app/myapp

# Menentukan port yang akan digunakan
EXPOSE 8080

# Menjalankan aplikasi
CMD ["./myapp"]