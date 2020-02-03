FROM alpine

WORKDIR /iwantoask

COPY static/ static/
COPY template/ template/
COPY iwantoask .

CMD ./iwantoask
