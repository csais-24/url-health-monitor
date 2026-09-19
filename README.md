# URL Health Monitor

## Features
- Check the health of a single URL
- Check the health of a list of URLs in a file, one URL per line
- Concurrent URL health checks
- Handles request timeout
- Reports HTTP status codes
- Docker support

## Usage

### Check a single URL
````bash
go run ./cmd/url-health-monitor --url https://httpbin.org
````
Checks the health of a single URL

### Check a list of URLs
````bash
go run ./cmd/url-health-monitor --url-list ./path-to-url-list.txt
````
Checks the health of a list of URLs in a file, one URL per line.

## Testing

### Checks for package functionality working correctly

```bash
go test -count=1 ./...
```
- Tests for a single url health check
- Check that HTTP status is the one expected
- Check that healthy and reachable are output correctly

### Check if there are any race condition
```bash
go test -count=1 -v -race ./...
```
- Checks that concurrency is working correctly and there are no race conditions.