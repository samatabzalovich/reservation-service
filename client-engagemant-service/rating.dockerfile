FROM alpine:latest
RUN apk --no-cache add curl

RUN mkdir /app

COPY ratingApp /app

CMD [ "/app/ratingApp"]