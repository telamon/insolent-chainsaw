#!/bin/bash
# this file is just a memo,
# committed for later
docker build -t wharfmaster_test -f DockerfileTest .
docker run --rm -it  -v $PWD:/go/src/github.com/telamon/wharfmaster -v /var/run/docker.sock:/tmp/docker.sock wharfmaster_test bash
