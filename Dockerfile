# Builder image
FROM golang:alpine AS builder
RUN apk --no-cache --virtual .build-dependencies add git curl openssh
WORKDIR $GOPATH/src/github.com/indiependente/gw-example
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /service

# Deployable image
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /service /app/
EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/service"]
