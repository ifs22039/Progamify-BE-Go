# Gunakan base image Go
FROM golang:1.21

# Set work directory
WORKDIR /app

# Copy project files
COPY . .

# Install dependencies & build
RUN go mod tidy
RUN go build -o main .

# Expose port yang digunakan
EXPOSE 8080

CMD [".cmd/api/main"]