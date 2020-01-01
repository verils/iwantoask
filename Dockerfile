FROM alpine

WORKDIR /app

COPY static/ static/
COPY template/ template/
COPY iwantoask .

CMD ./iwantoask