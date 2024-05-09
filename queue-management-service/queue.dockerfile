FROM alpine:latest

RUN mkdir /app

COPY queueApp /app

CMD [ "/app/queueApp"]