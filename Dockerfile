FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o support-portal ./cmd/server/

FROM scratch

COPY --from=builder /app/support-portal /support-portal
COPY /static ./static
COPY /templates ./templates

ENTRYPOINT ["/support-portal"]