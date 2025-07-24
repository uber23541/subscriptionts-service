FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o subscriptions ./server



FROM alpine:3.20
WORKDIR /root
ARG HTTP_PORT=8080
COPY --from=builder /app/subscriptions .
COPY .env .
EXPOSE ${HTTP_PORT}
CMD ["./subscriptions"]