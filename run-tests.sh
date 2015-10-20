#!/bin/bash
# this file is just a memo,
# committed for later
docker run --name=wharftest -it  -v $PWD:/app -v /var/run/docker.sock:/tmp/docker.sock wharfmaster_app test
