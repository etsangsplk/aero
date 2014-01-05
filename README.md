Aero
==========

This library provides backoff rules and rate limiting (similar to time.Ticker) to aid in making nice networked clients. Also, chainable http.RoundTrippers are provided for both to provide that functionality at the network level.

To **install** Aero, use `go get`:

    go get github.com/mceldeen/aero

This will then make the following packages available to you:

    github.com/mceldeen/aero/backoff
    github.com/mceldeen/aero/ratelimit


Backoff
----------

```go
// backoffer.go

package main

import (
    "fmt"
    "github.com/mceldeen/aero/backoff"
    "time"
)

var succeed = uint(8)

func do(n uint) bool {
    if n == succeed {
        return true
    }

    return false
}

func main() {
    b := backoff.NewExp()
    b.InitialDelay = 10 * time.Millisecond
    b.MaxDelay = 1000 * time.Millisecond
    b.FailAfter = 10

    n := uint(0)
    for delay, next := b.Next(n); next; delay, next = b.Next(n) {
        fmt.Printf("do request %d...", n)

        if do(n) {
            fmt.Println("success")
            return
        }

        fmt.Println("failed")

        fmt.Printf("wait %s...\n\n", delay)
        <-time.After(delay)
        n++
    }

    fmt.Println("too many requests")
}

```

```Shell
$ go run backoffer.go
do request 0...failed
wait 10ms...

do request 1...failed
wait 20ms...

do request 2...failed
wait 40ms...

do request 3...failed
wait 80ms...

do request 4...failed
wait 160ms...

do request 5...failed
wait 320ms...

do request 6...failed
wait 640ms...

do request 7...failed
wait 1s...

do request 8...success

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
