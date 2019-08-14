# Builder image
FROM golang:1.12-alpine AS builder
RUN apk --no-cache --virtual .build-dependencies add git curl openssh
ENV GO111MODULE=on
WORKDIR $GOPATH/src/github.com/indiependente/gw-example/
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /client ./client/

# Deployable image
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /client /app/
ENTRYPOINT ["/app/client"]