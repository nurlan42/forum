FROM golang:alpine

RUN mkdir /app /app/bin
RUN apk add build-base

LABEL author="Nurlan & Maksat"
LABEL description="forum"
LABEL Data="18.02.2022"


ADD . /app

WORKDIR /app

RUN go build -o bin/main ./cmd
RUN rm -rf cmd internal pkg server

CMD [ "/app/bin/main" ]

