# docker-wg
Docker WireGuard but with Http APIs

This project wraps WireGuard and some simple endpoints into one service.  
Those endpoints help you interact with WireGuard via Http requests,
when you run WireGuard as a VPN server in your microservice environment.


### Build
```shell
docker build . -t wg
```

### Run
```shell
docker run  \
  --cap-add=NET_ADMIN \
  --cap-add=SYS_MODULE \
  -e TZ=Europe/London \
  -e SERVERPORT=51820 \
  -e PEERS=1 \
  -e PEERDNS=auto \
  -e INTERNAL_SUBNET=10.13.13.0/24 \
  -e ALLOWEDIPS=0.0.0.0/0 \
  -p 51820:51820/udp \
  -p 8000:8000 \
  --name wg \
  wg
```

### API Usage
> Replace `wg0` with your interface name
1. Get list of peers
```shell
curl --request GET \
  --url http://localhost:8000/wg0/peer/list
```

2. Register new peer
```shell
curl --request POST \
  --url http://localhost:8000/wg0/peer/register
```

3. Remove a peer
```shell
curl --request DELETE \
  --url 'http://localhost:8000/wg0/peer?key=<Peer Public Key>
```

### References
1. https://github.com/WireGuard
2. https://hub.docker.com/r/linuxserver/wireguard
3. https://man7.org/linux/man-pages/man8/wg.8.html
4. https://man7.org/linux/man-pages/man8/wg-quick.8.html