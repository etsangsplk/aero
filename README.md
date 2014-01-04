Aero
==========

This library provides backoff rules and rate limiting (similar to time.Ticker) to aid in making nice networked clients. Also, chainable http.RoundTrippers are provided for both to provide that functionality at the network level.

Backoff
----------

```go

b := backoff.NewExp()
b.InitialDelay = 100 * time.Millisecond
b.MaxDelay = 500 * time.Millisecond
b.FailAfter = 10
b.JitterBefore = 0.01
b.JitterAfter = 0.2

// ...

var res *http.Response
var err error

n := uint(0)
for delay, next := b.Next(n); next; delay, next = b.Next(n) {
    res, err = do(req)

    if err == nil {
        return res, err
    }

    <-time.After(delay)
    n++
}

return res,err
```

Ratelimit
----------
```Go
// limiter.go

package main

import (
    "fmt"
    "github.com/mceldeen/aero/ratelimit"
    "time"
)

func main() {
    requests := make(chan int, 5)
    for i := 1; i <= 5; i++ {
        requests <- i
    }
    close(requests)

    // 1 req ever 200 milliseconds with a burst of 2
    limiter := ratelimit.NewBursty(1, 200*time.Millisecond, 2)

    for req := range requests {
        <-limiter.Tick()
        fmt.Println("request", req, time.Now())
    }
}
```

```Shell
$ go run limiter.go
request 1 2014-01-04 16:22:42.695883641 -0700 MST
request 2 2014-01-04 16:22:42.695965085 -0700 MST
request 3 2014-01-04 16:22:42.695979642 -0700 MST
request 4 2014-01-04 16:22:42.896342263 -0700 MST
request 5 2014-01-04 16:22:43.09632076 -0700 MST
```
