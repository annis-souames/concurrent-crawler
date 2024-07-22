## Distributed Web Crawler

This project is a simple concurrent web crawler to learn concurrency concepts in Go such as Goroutines, channels, mutexes, locks, and wait groups.

To run this project, simply clone the project, install the `net/html` package, build it and then run it:

```
go get golang.org/x/net/html
go build .
go run .
```

The crawler should fetch 200 links aynchronously, if executed in a multi-core CPU, the crawler will perform several crawling operations in parallel. This crawler was tested on a 6-core CPU.


**Find a full tutorial about this crawler on my blog: https://souames.hashnode.dev/learning-concurrency-building-a-concurrent-web-crawler-with-go **
