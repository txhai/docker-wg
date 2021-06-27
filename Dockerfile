FROM ghcr.io/linuxserver/wireguard
MAINTAINER 'https://github.com/txhai'
RUN \
    echo "**** install dependencies ****" && \
    apt-get update && \
    apt-get install -y --no-install-recommends wget && \
    wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
    mkdir /root/.conda && \
    bash Miniconda3-latest-Linux-x86_64.sh -b && \
    rm -f Miniconda3-latest-Linux-x86_64.sh && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/*

ADD ./src /http

RUN cd /http && \
    /root/miniconda3/bin/python -m pip install -r requirements.txt

CMD ["/root/miniconda3/bin/python", "/http/entry.py"]