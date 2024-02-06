FROM alpine:latest

RUN mkdir /app

COPY staffApp /app

CMD [ "/app/staffApp"]
