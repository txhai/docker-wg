FROM ghcr.io/linuxserver/wireguard
MAINTAINER 'https://github.com/txhai'

WORKDIR /http

ADD ./src .

RUN curl -k -o "conda.sh" "https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh" && \
    mkdir /root/.conda && \
    bash conda.sh -b && \
    /root/miniconda3/bin/python -m pip install -r ./requirements.txt && \
    rm -f conda.sh && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/*

CMD ["/root/miniconda3/bin/python", "-u", "/http/entry.py"]