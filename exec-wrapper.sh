#!/bin/bash
cmd="docker exec wharfmaster_app_1 $@"
echo "Wrap: $cmd"
$cmd