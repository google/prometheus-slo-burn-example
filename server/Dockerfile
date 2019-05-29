FROM golang:1.12-alpine as builder
ENV GO111MODULE=on
WORKDIR /go/src/example
RUN apk add --no-cache git
COPY . .

RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/helloworld

ENV PORT 8080
EXPOSE $PORT

CMD ["/go/bin/helloworld"]
