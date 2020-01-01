FROM alpine

WORKDIR /app

COPY static/ statistic/
COPY template/ template/
COPY iwantoask .

CMD ./iwantoask