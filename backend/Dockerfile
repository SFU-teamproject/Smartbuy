FROM golang:1.24.1-alpine AS builder
WORKDIR /smartbuy
# install git with no caching (to reduce space)
RUN apk add --no-cache git
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=linux go build -o /tmp/smartbuy

# Final stage
FROM scratch
WORKDIR /
COPY --from=builder /tmp/smartbuy .
EXPOSE 8081
CMD ["/smartbuy"]