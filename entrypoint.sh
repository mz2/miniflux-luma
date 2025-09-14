#!/bin/sh
set -e

if [ -n "$MINIFLUX_API_TOKEN" ]; then
  echo "$MINIFLUX_API_TOKEN" > /tmp/api_token

  # Write feed format to file if specified
  FORMAT_ARG=""
  if [ -n "$FEED_FORMAT" ]; then
    echo "$FEED_FORMAT" > /tmp/feed_format
    FORMAT_ARG="-format-file /tmp/feed_format"
  fi

  exec ./miniflux-luma -api-token-file /tmp/api_token $FORMAT_ARG "$@"
else
  echo "Error: MINIFLUX_API_TOKEN environment variable not set"
  echo "Please add MINIFLUX_API_TOKEN to your .env file"
  exit 1
fi