# miniflux-luma

Atom/RSS feed exporter for Miniflux starred items

## Installation

```
go get -u github.com/erdnaxe/miniflux-luma
```

## Usage

Fetch your API token from Miniflux settings and pass it via the `-api-token` argument.

Then you may start the web service:

```
miniflux-luma -api-token YOUR_TOKEN -endpoint https://rss.example.com -listen-addr 127.0.0.1:8080
```

By default, the 10 most recent starred items are included. To change this:

```
miniflux-luma -api-token YOUR_TOKEN -limit 0   # Include ALL starred items
miniflux-luma -api-token YOUR_TOKEN -limit 50  # Include 50 most recent starred items
```

## Feed Formats

The service supports both Atom (default) and RSS formats:

### Endpoints
- `/` - Returns feed in the configured default format
- `/atom` - Always returns Atom format
- `/rss` - Always returns RSS format

### Configuration

To set the default format, use the `-format` argument:

```
miniflux-luma -api-token YOUR_TOKEN -format rss -endpoint https://rss.example.com
```

