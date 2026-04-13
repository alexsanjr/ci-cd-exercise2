# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api-go .


FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /out/api-go /usr/local/bin/api-go

ENV PORT=3000
EXPOSE 3000

USER app

CMD ["api-go"]
