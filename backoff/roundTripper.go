package backoff

import (
    "net/http"
    "time"
)

type HttpTransport struct {
    Transport http.RoundTripper

    b   Backoffer
}

func NewHttpTransport(b Backoffer) *HttpTransport {
    return &HttpTransport{
        Transport: http.DefaultTransport,
        b:         b,
    }
}

func (t *HttpTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
    n := uint(0)

    for delay, next := t.b.Next(n); next; delay, next = t.b.Next(n) {
        res, err = t.Transport.RoundTrip(req)

        if err == nil {
            return
        }

        <-time.After(delay)
        n++
    }

    return
}
