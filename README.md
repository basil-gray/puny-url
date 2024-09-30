# PunyURL 

A lightweight URL shortener written in Go that utilises: 

1. SQLite for persistent storage
2. The concurrency-safe sync.Map data structure to handle concurrent requests and cache URLs that see frequent access

To improve scalability, each URL has a TTL value configurable by passing an arguement. The TTL is refreshed at every successful cache hit, or whenever a trip is made to the database. URLs that have exceeded their TTL are evicted from the cache by a job running in a goroutine, also at a configurable interval.

When attempting to generate a short URL for a duplicate long one, PunyURL will attempt to first find it in the cache, and then the database. 

### Screenshot
<img width="1215" alt="image" src="https://github.com/user-attachments/assets/7b2cbe19-d230-45f5-bc75-ad85d66ad2bc">

# Usage
## Docker
Run docker compose to build and run the container. To test that it's working, you can run the given curl script.

```
git clone https://github.com/basil-gray/puny-url.git

docker-compose up --build

curl --location 'http://localhost:8080/shorten' \
--header 'Content-Type: application/json' \
--data '{"long_url": "http://google.com"}'
```
## Locally
Alternatively, you can build the project and run it from terminal, or simply run main.go.

```
// If you don't have sqlite3 installed:
brew install sqlite

go build
./puny-url

go run .
```

## Flags

```
// Silent flag to prevent logging cache hits/misses
-s

// The interval of the eviction job (default: 1s)
-e 30s

// The TTL for a URL in the cache (default: 30s)
-ttl 30s

-port 8080

These flags are DurationVar, so 30m and 1h are also accepted values.

```

# Tests

There are currently only tests that check for accepting and rejecting valid and invalid URLs. You can run them like so:

```go test ./tests ```