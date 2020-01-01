FROM golang:1.13.5-stretch

WORKDIR /app

COPY iwantoask .

CMD ./iwantoask