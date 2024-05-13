FROM alpine:latest
RUN apk --no-cache add curl
RUN mkdir /app

COPY staffApp /app

CMD [ "/app/staffApp"]
