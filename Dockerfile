FROM armv7/armhf-debian


RUN apt-get update && \
    apt-get install -y \
    build-essential wget \
    libpcre3-dev zlib1g-dev libssl-dev

ENV HA_VER_MAJOR=1.6 \
    HA_VER_MINOR=0 \
    HA_OPTS="TARGET=linux2628 USE_PCRE=1 USE_OPENSSL=1 USE_ZLIB=1"

RUN cd /usr/local/src/ && \
    wget -O haproxy.tar.gz http://www.haproxy.org/download/${HA_VER_MAJOR}/src/haproxy-${HA_VER_MAJOR}.${HA_VER_MINOR}.tar.gz && \
    tar xfvz haproxy.tar.gz && \
    cd haproxy-${HA_VER_MAJOR}.${HA_VER_MINOR} && \
    make -j4 ${HA_OPTS} && \
    make install || echo "Done setting up HA-Proxy"

VOLUME /app
WORKDIR /app
ENV HA_CONF=/app/haproxy.conf
CMD haproxy -db -f $HA_CONF

