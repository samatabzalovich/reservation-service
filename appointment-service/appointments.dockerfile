FROM alpine:latest
RUN apk --no-cache add curl
RUN mkdir /app

COPY appointmentApp /app

CMD [ "/app/appointmentApp"]


