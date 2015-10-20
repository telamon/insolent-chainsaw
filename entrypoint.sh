#!/bin/bash 
set -e
# Purge / initialize vhosts file
echo "" > /app/vhosts.conf
# Let nginx fork off to background.
exec nginx & 
exec "$@"