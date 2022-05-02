FROM golang:1.18 as builder
WORKDIR /usr/src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo -o service .

FROM alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
WORKDIR /usr/app
COPY --from=builder /usr/src/service .
COPY static .
COPY tpl .
CMD ["service"]
