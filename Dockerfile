FROM golang:1.26rc1-alpine3.23

COPY . .

RUN go mod download

ENTRYPOINT ["go", "run"]

CMD ["./cmd/url-shortener/url-shortener.go"]