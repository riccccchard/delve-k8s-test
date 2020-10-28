FROM golang:1.14 as builder

WORKDIR /build

COPY . .

ENV GO111MODULE=on \
        CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64

RUN go build -o delveServer

WORKDIR /dist

RUN cp /build/delveServer .

EXPOSE 3333

CMD ["/dist/delveServer"]