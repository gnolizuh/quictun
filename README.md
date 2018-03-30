# quictun

The simplest tunnel service based on QUIC.

### QuickStart

Download precompiled [Releases](https://github.com/gnolizuh/quictun/releases).

```
QUIC Client: ./quictun-client-darwin-amd64 -l ":1935" -r "QUIC_SERVER_IP:6935"
QUIC Server: ./quictun-server-darwin-amd64 -l ":6935" -t "TARGET_IP:1935"
```

The above commands will establish port forwarding for 1935/tcp as:

> Application -> **QUIC Client(1935/tcp) -> QUIC Server(6935/udp)** -> Target Server(1935/tcp) 

Tunnels the original connection:

> Application -> Target Server(1935/tcp) 

### Install from source

```
$go get -u github.com/gnolizuh/quictun/client
$go get -u github.com/gnolizuh/quictun/server
```