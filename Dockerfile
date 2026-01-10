FROM golang:1.26rc1-alpine3.23 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager ./cmd/url-shortener/url-shortener.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /workspace/manager .

USER nonroot

ENTRYPOINT ["/app/manager"]
