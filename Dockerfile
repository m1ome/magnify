FROM golang:1.23-alpine3.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o /magnify

FROM alpine:3.20

WORKDIR /

COPY --from=build-stage /magnify /magnify

EXPOSE 8080

ENTRYPOINT ["/magnify"]
