FROM golang:1.21
WORKDIR /opt/app
COPY socks5-relay /opt/app/
CMD ["./socks5-relay"]
