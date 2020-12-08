FROM golang:alpine as go-builder

RUN mkdir /photomailer
WORKDIR /photomailer
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o /go/bin/photomailer

FROM alpine

WORKDIR /go/bin
COPY --from=go-builder /go/bin/photomailer .
COPY .env .
EXPOSE 80
ENTRYPOINT ["./photomailer"]
