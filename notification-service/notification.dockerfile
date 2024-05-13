FROM alpine:latest
RUN apk --no-cache add curl
RUN mkdir /app

COPY notificationApp /app

CMD [ "/app/notificationApp"]