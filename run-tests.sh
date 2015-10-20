#!/bin/bash
# this file is just a memo,
# committed for later
name=wharfmaster_test
docker run --name=$name -it  -v $PWD:/app -v /var/run/docker.sock:/tmp/docker.sock wharfmaster_app test || docker start -ai $name
