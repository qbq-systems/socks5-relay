# Description

Creates a connections with a socks5 proxies using authorization. A local connections is available on a random ports and without authorization.

The main purpose is to use with those clients that do not support socks5 authorization, for example Chrome, Puppeteer and so on.

# Install

### via `go build`

```
git clone https://github.com/qbq-systems/socks5-relay
cd socks5-relay
go build
#
# to run a REST server on port 18000
#
./socks5-relay 
#
# or custom port
#
LISTEN_PORT=19876 ./socks5-relay
```

### via `Docker`

```
docker pull qbqsystemsbot/socks5-relay:latest
docker run --name qbq-socks5-relay qbqsystemsbot/socks5-relay:latest
```

# REST API

### create relay

```
curl -XPOST \
  -H "Content-Type: application/json" \
  -d '{"host": "1.2.3.4", "port": "12345", "user": "user", "password": "password"}' http://127.0.0.1:18000
{"port":38153,"message":"Added"}
```

### use relay

```
curl -XGET -x 127.0.0.1:38153 https://httpbin.org/get
{
  "args": {}, 
  "headers": {
    "Accept": "*/*", 
    "Host": "httpbin.org", 
    "User-Agent": "curl/7.74.0", 
    "X-Amzn-Trace-Id": "Root=1-65461b1c-664646b277b9482207aae14a"
  }, 
  "origin": "1.2.3.4", 
  "url": "https://httpbin.org/get"
}
```

### delete relay

```
curl -XDELETE \
  -H "Content-Type: application/json" \
  -d '{"host": "1.2.3.4", "port": "12345"}' http://127.0.0.1:18000
{"port":0,"message":"Deleted"}
```