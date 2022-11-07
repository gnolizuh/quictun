# quictun

The simplest tunnel service based on QUIC.

## Benchmarks

-   **sshping**

```bash
                            QUIC            TCP                            
ssh-Login-Time:              2.47  s         2.60  s
Minimum-Latency:             58.3 ms         59.0 ms
Median-Latency:              59.5 ms         59.8 ms
Average-Latency:             60.1 ms         60.1 ms
Average-Deviation:           3.32 ms         2.15 ms
Maximum-Latency:              105 ms         96.2 ms
Echo-Count:                 1.00 kB         1.00 kB
Upload-Size:                8.00 MB         8.00 MB
Upload-Rate:                 811 kB/s        810 kB/s
Download-Size:              8.00 MB         8.00 MB
Download-Rate:              2.45 MB/s       1.80 MB/s
```

-   **iperf3**
### QUIC

```bash
[ ID] Interval           Transfer     Bitrate         Retr  Cwnd
[  5]   0.00-1.00   sec  11.2 MBytes  94.3 Mbits/sec    2   1.69 MBytes       
[  5]   1.00-2.00   sec  1.25 MBytes  10.5 Mbits/sec    0    639 KBytes       
[  5]   2.00-3.00   sec  0.00 Bytes  0.00 bits/sec    0    639 KBytes       
[  5]   3.00-4.00   sec  1.25 MBytes  10.5 Mbits/sec    0    639 KBytes       
[  5]   4.00-5.00   sec  1.25 MBytes  10.5 Mbits/sec    0    639 KBytes       
[  5]   5.00-6.00   sec  0.00 Bytes  0.00 bits/sec    0    639 KBytes       
[  5]   6.00-7.00   sec  1.25 MBytes  10.5 Mbits/sec    0    639 KBytes       
[  5]   7.00-8.00   sec  1.25 MBytes  10.5 Mbits/sec    0    639 KBytes       
[  5]   8.00-9.00   sec  0.00 Bytes  0.00 bits/sec    0    639 KBytes       
[  5]   9.00-10.00  sec  1.25 MBytes  10.5 Mbits/sec    0    639 KBytes       
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate         Retr
[  5]   0.00-10.00  sec  18.8 MBytes  15.7 Mbits/sec    2             sender
[  5]   0.00-12.29  sec  9.69 MBytes  6.61 Mbits/sec                  receive
```
### TCP
```bash
[ ID] Interval           Transfer     Bitrate         Retr  Cwnd
[  5]   0.00-1.00   sec  10.0 MBytes  83.8 Mbits/sec    1   5.31 MBytes       
[  5]   1.00-2.00   sec  1.25 MBytes  10.5 Mbits/sec    0   5.31 MBytes       
[  5]   2.00-3.00   sec  1.25 MBytes  10.5 Mbits/sec    2   5.31 MBytes       
[  5]   3.00-4.00   sec  0.00 Bytes  0.00 bits/sec    2   5.31 MBytes       
[  5]   4.00-5.00   sec  1.25 MBytes  10.5 Mbits/sec    2   5.31 MBytes       
[  5]   5.00-6.00   sec  1.25 MBytes  10.5 Mbits/sec    2   5.31 MBytes       
[  5]   6.00-7.00   sec  0.00 Bytes  0.00 bits/sec    3   5.31 MBytes       
[  5]   7.00-8.00   sec  1.25 MBytes  10.5 Mbits/sec    2   5.31 MBytes       
[  5]   8.00-9.00   sec  1.25 MBytes  10.5 Mbits/sec    2   2.62 MBytes       
[  5]   9.00-10.00  sec  0.00 Bytes  0.00 bits/sec    1   1.31 MBytes       
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate         Retr
[  5]   0.00-10.00  sec  17.5 MBytes  14.7 Mbits/sec   17             sender
[  5]   0.00-10.57  sec  8.75 MBytes  6.95 Mbits/sec                  receiver
```


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