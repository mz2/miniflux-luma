#!/bin/sh
set -e

if [ -n "$MINIFLUX_API_TOKEN" ]; then
  echo "$MINIFLUX_API_TOKEN" > /tmp/api_token
  exec ./miniflux-luma -api-token-file /tmp/api_token "$@"
else
  echo "Error: MINIFLUX_API_TOKEN environment variable not set"
  echo "Please add MINIFLUX_API_TOKEN to your .env file"
  exit 1
fi