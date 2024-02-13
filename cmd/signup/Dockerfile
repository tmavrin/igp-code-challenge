# build stage
FROM golang:1.21-alpine AS build-env

RUN apk --no-cache add git
ENV GOPRIVATE=github.com/tmavrin

WORKDIR /go/src/github.com/tmavrin/igp-code-challenge/
COPY . .
RUN go build -o service cmd/api/main.go

# final image
FROM alpine:3.19
RUN apk update && apk add ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.com/tmavrin/igp-code-challenge/service /app/
ENTRYPOINT /app/service
