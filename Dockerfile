FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/drcom_go .


FROM alpine:latest
RUN apk add --no-cache libc6-compat
WORKDIR /drcom
COPY --from=builder /app/drcom_go .
CMD ["./drcom_go","--username","zhangs2121","--password" ,"Pa$$w0rb","--mac","12:34:56:78:9a:bc"]