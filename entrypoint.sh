#!/bin/sh
set -e

if [ -n "$MINIFLUX_API_TOKEN" ]; then
  echo "$MINIFLUX_API_TOKEN" > /tmp/api_token

  # Set feed format from environment variable, default to atom
  FORMAT_ARG=""
  if [ -n "$FEED_FORMAT" ]; then
    FORMAT_ARG="-format $FEED_FORMAT"
  fi

  exec ./miniflux-luma -api-token-file /tmp/api_token $FORMAT_ARG "$@"
else
  echo "Error: MINIFLUX_API_TOKEN environment variable not set"
  echo "Please add MINIFLUX_API_TOKEN to your .env file"
  exit 1
fi