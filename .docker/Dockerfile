FROM golang:1.21-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /conductorProxy ./cmd/main.go

FROM alpine:3.19

WORKDIR /

COPY --from=build-stage /conductorProxy /conductorProxy

EXPOSE 8000

ENTRYPOINT [ "/conductorProxy" ]