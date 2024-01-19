**cache package**

Package `cache` provides caching functionality for a service satisfying `cache.Servicer` interface. 

**Usage**
```golang
userId := "user_1"
service := &UserService{ // UserService satisfies `cache.Servicer` interface
    // GetOne(id string)
} 

cachedService := cache.NewCache(service)

user, err := cachedService.GetOne(userId)
// handle `err` and use `user`
```

**Testing**
```
go test ./...
```

**Considered approaches**

`sync.Mutex`
- it makes the cache slower as each request is handled synchronuously due to cache struct being locked on each read/write
- would be suitable for write heavy usecase

`sync.RWMutex`
- allows asynchronous reading, however until entry is succesfully inserted into the cache multiple concurrent request might miss the cache simultaneusly
- synchronously warming up the cache allows for predictable db hits value in tests

`sync.Map`
- similar behaviour to `sync.RWMutex`, optimised under the hood at the cost of strict typing

Decided to use `sync.RWMutex` for type safety.
