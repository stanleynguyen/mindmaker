FROM golang:1.13-alpine as builder
WORKDIR /go/src/github.com/stanleynguyen/mindmaker
RUN apk update && apk upgrade
COPY . .
RUN GOOS=linux go build -o mindmaker.out .

FROM alpine:latest
RUN apk update && apk upgrade
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/stanleynguyen/mindmaker/mindmaker.out .
CMD ["./mindmaker.out"]
