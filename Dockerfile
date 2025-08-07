# --- Build stage ---
FROM golang:1.24.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o notifier


# --- Run stage ---
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /app/notifier .
COPY --from=builder /app/store/templates /store/templates

ENTRYPOINT ["/notifier"]
