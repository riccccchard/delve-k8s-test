FROM golang:1.14 as builder

WORKDIR /build

COPY . .

ENV GO111MODULE=on \
        CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64

RUN go build -mod=vendor -o delveClient

WORKDIR /dist

RUN cp /build/delveClient .

EXPOSE 8888

CMD ["/dist/delveClient"]