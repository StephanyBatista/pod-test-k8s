FROM golang:1.20

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /decode-jwt

EXPOSE 3015

# Run
CMD ["/decode-jwt"]