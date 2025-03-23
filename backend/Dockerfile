FROM golang:1.24.1-alpine AS builder
WORKDIR /smartbuy
# install git with no caching (to reduce space)
RUN apk add --no-cache git
# Copy dependency files first
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy the rest of the code
COPY . .
# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smartbuy

# Final stage
FROM scratch
WORKDIR /
COPY --from=builder /smartbuy/smartbuy /smartbuy
EXPOSE 8081
CMD ["/smartbuy"]