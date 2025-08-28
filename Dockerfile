FROM docker.io/golang:1.25.0-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -o bin/kepegawaian ./services/kepegawaian
RUN go build -o bin/portal ./services/portal

FROM docker.io/alpine:3.22

WORKDIR /app

COPY --from=builder /app/bin/* .

EXPOSE 8000
