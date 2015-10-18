FROM debian
#FROM armv7/armhf-debian

MAINTAINER Tony Ivanov "telamohn@gmail.com"

ENV NGINX_VERSION=1.9.5 \
    NGINX_OPTS="--with-http_ssl_module --with-pcre-jit --with-http_stub_status_module" \
    MAKE_OPTS="-j8"

RUN deps='curl ca-certificates'; \
    builddeps='build-essential libpcre3-dev zlib1g-dev libssl-dev'; \
    set -x; \
    apt-get update && \
    apt-get install -y --no-install-recommends $deps $builddeps && \
    cd /usr/local/src && \
    curl -o nginx.tar.gz http://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz && \
    tar xfvz nginx.tar.gz && \
    rm nginx.tar.gz && \
    mv nginx-${NGINX_VERSION} nginx && \
    cd nginx && \
    ./configure --prefix=/opt/nginx \
        --sbin-path=/usr/sbin/nginx \
        $NGINX_OPTS && \ 
    make $MAKE_OPTS && \
    make install && \
    cd .. && rm -rf nginx && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false -o APT::AutoRemove::SuggestsImportant=false $builddeps && \
    apt-get clean -y
    #rm -rf /var/lib/apt/lists/*

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /opt/nginx/logs/access.log && \
    ln -sf /dev/stderr /opt/nginx/logs/error.log

# install docker-gen
ENV DGVERSION 0.4.2
ENV DOWNLOAD_URL https://github.com/jwilder/docker-gen/releases/download/$DGVERSION/docker-gen-linux-amd64-$DGVERSION.tar.gz
ENV DOCKER_HOST unix:///tmp/docker.sock
RUN curl -o docker-gen.tar.gz -L $DOWNLOAD_URL \
    && tar -C /usr/local/bin -xvzf docker-gen.tar.gz \
    && rm docker-gen.tar.gz

EXPOSE 80 443

WORKDIR /app
RUN mkdir -p /app/certs

ENV NGINX_CONF="/app/nginx.conf"

CMD ["nginx", "-g", "-f $NGINX_CONF" ,"daemon off;"]