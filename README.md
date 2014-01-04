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
