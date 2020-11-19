FROM golang:1.15.5-alpine3.12

RUN mkdir /build

ADD . /build

WORKDIR /build

RUN go build -o main .

CMD ["/build/main"]
