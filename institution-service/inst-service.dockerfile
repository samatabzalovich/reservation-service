FROM alpine:latest
RUN apk --no-cache add curl
RUN mkdir /app

COPY instApp /app

CMD [ "/app/instApp"]