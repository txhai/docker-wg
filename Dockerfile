FROM ghcr.io/linuxserver/baseimage-ubuntu:jammy AS build
# Build
RUN curl -s https://go.dev/dl/go1.18.3.linux-amd64.tar.gz -L | tar -C /usr/local -xz
ENV PATH $PATH:/usr/local/go/bin
ADD . /wgapp
WORKDIR /wgapp
RUN mkdir /wgapp/bin && GOOS=${GOOS} GOARCH=${GOARCH} go build -v -o /wgapp/bin/wgapi /wgapp/main.go

FROM ghcr.io/txhai/docker-wireguard:latest
COPY --from=build /wgapp/bin/ /usr/bin/
CMD ["/usr/bin/with-contenv", "/usr/bin/wgapi"]