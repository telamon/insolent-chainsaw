#!/bin/bash
docker build -f Dockerfile-1.0.0 -t telamon/wharftest-version:1.0.0 .
docker build -f Dockerfile-1.1.0 -t telamon/wharftest-version:1.1.0 .
docker build -f Dockerfile-1.2.0 -t telamon/wharftest-version:1.2.0 .
docker tag telamon/wharftest-version:1.2.0 telamon/wharftest-version:latest
