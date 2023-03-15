# FROM mongo:6.0.4
FROM mongo:latest

RUN openssl rand -base64 756 > /keyfile
RUN chmod 400 /keyfile
RUN chown 999:999 /keyfile

CMD ["mongod", "--replSet", "rs0", "--bind_ip_all", "--keyFile", "/keyfile"]
