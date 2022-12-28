package mid

// Throttle per client or endpoint by adding that reference (IP or r.URL.Path)
// as key in the short-lived (TTL) cache;
// refusing the request if its reference exists in cache.
// TTLCache : https://github.com/ReneKroon/ttlcache (thread safe)

// See @ /DEV/go/labs/14-throttle/main.go
