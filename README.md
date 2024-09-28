# PunyURL 

A lightweight URL shortener written in Go that utilises: 

1. SQLite for persistent storage
2. The concurrency-safe sync.Map data structure to handle concurrent requests and cache URLs that see frequent access

To improve scalability, each URL has a TTL value configurable via an arguement. The TTL value is refreshed at every successful cache hit. URLs that have exceeded their TTL are evicted from the cache by a job running in a goroutine, also at a configurable interval.

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
-- The eviction interval (default: 30s)
-e 30s

-- The TTL for a URL in the cache (default: 30s)
-ttl 30s

These flags are DurationVar, so 30m and 1h are also accepted values.

```

# Tests

There are currently only tests that check for valid and invalid URLs. You can run them like so:

```go test ./... ```
