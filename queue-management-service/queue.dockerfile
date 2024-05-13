FROM alpine:latest
RUN apk --no-cache add curl
RUN mkdir /app

COPY queueApp /app

CMD [ "/app/queueApp"]