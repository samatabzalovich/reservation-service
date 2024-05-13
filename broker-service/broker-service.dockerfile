
FROM alpine:latest

RUN apk --no-cache add curl

RUN mkdir /app

COPY brokerApp /app

EXPOSE 8080

CMD [ "/app/brokerApp"]