FROM alpine:latest

WORKDIR /app

COPY qwq-server .

RUN chmod +x ./qwq-server

EXPOSE 8080

ENTRYPOINT ["./qwq-server"]