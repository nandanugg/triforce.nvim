ARG CI_REGISTRY=docker.io
FROM $CI_REGISTRY/golang:1.25.0-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/kepegawaian ./services/kepegawaian
RUN go build -o bin/portal ./services/portal

FROM $CI_REGISTRY/alpine:3.22

ARG CA_CRT=LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
RUN echo $CA_CRT | base64 -d >> /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

COPY --from=builder /app/bin/* .

EXPOSE 8000
