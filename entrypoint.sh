#!/bin/bash 
set -e
# Purge / initialize vhosts file
echo "" > /app/vhosts.conf
# Let nginx fork off to background.
exec nginx & 

if [ "$1" = "test" ]; then
	export GOLANG_VERSION="1.5.1"
	export GOLANG_DOWNLOAD_URL="https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz"
	export GOLANG_DOWNLOAD_SHA1="46eecd290d8803887dec718c691cc243f2175fe0"
	export GOPATH="/go"
	export PATH="$GOPATH/bin:/usr/local/go/bin:$PATH"
	export PROJECT_PATH="/go/src/github.com/telamon/wharfmaster"

	if [ ! -f "/app/wharfmaster.go" ]; then
		echo "Test command is intended to be used with the source-code volumed on /app -v $PWD:/app"
		exit 1
	fi
	# Initialize destenvironment
	if [ ! -L $PROJECT_PATH ]; then
		if [ ! -f "golang.tar.gz" ]; then
			echo "Setting up test-environment" \
			&& curl -fL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
			&& echo "$GOLANG_DOWNLOAD_SHA1  golang.tar.gz" | sha1sum -c - 
		fi
		if [ ! -d "$GOPATH" ]; then
			tar -C /usr/local -xzf golang.tar.gz \
			&& mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" \
			&& apt-get install -y git \
			&& go get github.com/onsi/ginkgo/ginkgo \
			&& go get github.com/onsi/gomega 
		fi
		mkdir -p $PROJECT_PATH  && rm -r $PROJECT_PATH \
		&& ln -sf /app $PROJECT_PATH 
	fi
	cd $PROJECT_PATH && go get && ginkgo watch
else
	exec "$@"
fi
