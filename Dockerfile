FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o energy-monitor .

FROM alpine:latest

RUN apk add --no-cache sqlite-libs tzdata

WORKDIR /app

COPY --from=builder /app/energy-monitor .
COPY static ./static

RUN mkdir -p /app/data

ENV TZ=Europe/Stockholm
ENV DB_PATH=/app/data/energy.db
ENV PORT=8081

EXPOSE 8081

CMD ["./energy-monitor"]
