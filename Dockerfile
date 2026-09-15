ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
# Retry to ride out transient Go module proxy errors
RUN for i in 1 2 3; do go mod download && break || sleep 5; done
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o baton-appstoreconnect ./cmd/baton-appstoreconnect

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /app/baton-appstoreconnect /baton-appstoreconnect
ENTRYPOINT ["/baton-appstoreconnect"]
