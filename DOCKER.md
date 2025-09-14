# Docker Configuration

## Environment Variables

- `MINIFLUX_API_TOKEN` (required): Your Miniflux API token
- `FEED_FORMAT` (optional): Feed format - `atom` (default) or `rss`

## Endpoints

The service provides the following endpoints:
- `/` - Returns feed in the default format (configured via FEED_FORMAT)
- `/atom` - Always returns Atom format
- `/rss` - Always returns RSS format

## Example

```bash
docker run -e MINIFLUX_API_TOKEN=your_token -e FEED_FORMAT=rss miniflux-luma
```