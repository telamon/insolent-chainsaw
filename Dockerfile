FROM debian
#FROM armv7/armhf-debian

MAINTAINER Tony Ivanov "telamohn@gmail.com"

RUN apt-get update && \
    apt-get install -y \
    build-essential wget \
    libpcre3-dev zlib1g-dev libssl-dev ca-certificates

ENV NGINX_VERSION=1.9.5 \
    NGINX_OPTS="--with-http_ssl_module --with-pcre-jit --with-http_stub_status_module" \
    MAKE_OPTS="-j8"

RUN cd /usr/local/src && \
    wget -O nginx.tar.gz http://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz && \
    tar xfvz nginx.tar.gz && \
    rm nginx.tar.gz && \
    mv nginx-${NGINX_VERSION} nginx && \
    cd nginx && \
    ./configure --prefix=/opt/nginx \
        --sbin-path=/usr/sbin/nginx \
        $NGINX_OPTS && \ 
    make $MAKE_OPTS && \
    make install

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /opt/nginx/logs/access.log
RUN ln -sf /dev/stderr /opt/nginx/logs/error.log

EXPOSE 80 443

VOLUME /app
WORKDIR /app

ENV NGINX_CONF="/app/nginx.conf"

CMD ["nginx", "-g", "-f $NGINX_CONF" ,"daemon off;"]