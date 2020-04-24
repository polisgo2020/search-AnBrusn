FROM golang:1.14 as builder

WORKDIR /usr/src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -installsuffix cgo -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /usr/src
COPY --from=builder /usr/src/app .
CMD ./app search -http
