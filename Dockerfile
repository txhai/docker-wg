FROM ghcr.io/linuxserver/baseimage-ubuntu:focal AS build
# Build
# RUN apk add build-base --update
RUN curl -s https://go.dev/dl/go1.18.3.linux-amd64.tar.gz -L | tar -C /usr/local -xz
ENV PATH $PATH:/usr/local/go/bin
ADD . /app
WORKDIR /app
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -v -o /app/bin/api /app/main.go

FROM ghcr.io/linuxserver/wireguard
COPY --from=build /app/bin/ /bin/
CMD ["api"]