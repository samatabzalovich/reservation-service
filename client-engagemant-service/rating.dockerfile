FROM alpine:latest

RUN mkdir /app

COPY appointmentApp /app

CMD [ "/app/appointmentApp"]