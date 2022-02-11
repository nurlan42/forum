FROM golang:1.16.3

RUN mkdir /app

LABEL author="Nurlan & Maksat"
LABEL description="forum"


ADD . /app

WORKDIR /app

RUN go build -o main .

CMD [ "app/main" ]

