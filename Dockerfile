FROM golang:alpine

RUN mkdir /app
RUN apk add build-base

LABEL author="Nurlan & Maksat"
LABEL description="forum"
LABEL Data="18.02.2022"


ADD . /app

WORKDIR /app

RUN go build -o main ./cmd

CMD [ "/app/main" ]

