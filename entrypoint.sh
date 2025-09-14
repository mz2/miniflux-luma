#!/bin/sh
set -e

if [ -z "$MINIFLUX_API_TOKEN" ]; then
  echo "Error: MINIFLUX_API_TOKEN environment variable not set"
  echo "Please add MINIFLUX_API_TOKEN to your .env file"
  exit 1
fi

# Set feed format from environment variable, default to atom
FORMAT="${FEED_FORMAT:-atom}"

exec ./miniflux-luma -api-token "$MINIFLUX_API_TOKEN" -format "$FORMAT" "$@"