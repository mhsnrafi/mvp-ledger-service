FROM golang:1.19-bullseye

RUN apt update

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Change .env with docker
COPY .env.local ./.env
COPY ../nginx.conf /etc/nginx/nginx.conf


# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]