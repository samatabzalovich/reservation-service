FROM alpine:latest

RUN mkdir /app

COPY instApp /app

CMD [ "/app/instApp"]