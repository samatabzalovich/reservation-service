FROM alpine:latest
# Install curl
RUN apk --no-cache add curl
RUN mkdir /app

COPY authApp /app

EXPOSE 8090

CMD [ "/app/authApp"]