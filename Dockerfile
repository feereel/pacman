#Building stage
FROM golang:1.22.3-alpine3.20 AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

#Container
FROM alpine:3.20

WORKDIR /app/bin

COPY --from=builder /app/bin/main .

EXPOSE  4444

CMD ["/app/bin/main"]
