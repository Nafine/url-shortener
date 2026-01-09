FROM golang:1.26rc1-alpine3.23

COPY . .

RUN addgroup -g 1001 shortener && \
    adduser -D -u 1001 -G shortener shortener

RUN go mod download

ENTRYPOINT ["go", "run"]

CMD ["./cmd/url-shortener/url-shortener.go"]