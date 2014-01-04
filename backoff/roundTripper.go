package backoff

import (
    "errors"
    "net/http"
    "os"
    "os/signal"
    "time"
)

type HttpRoundTripper struct {
    Transport http.RoundTripper

    b   Backoffer
}

func NewHttpRoundTripper(b Backoffer) *HttpRoundTripper {
    return &HttpRoundTripper{
        Transport: http.DefaultTransport,
        b:         b,
    }
}

type rtres struct {
    res *http.Response
    err error
}

func (t *HttpRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
    n := uint(0)

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    defer signal.Stop(c)

    rtrc := make(chan *rtres)

    for delay, next := t.b.Next(n); next; delay, next = t.b.Next(n) {
        go func() {
            rtr := new(rtres)
            rtr.res, rtr.err = t.Transport.RoundTrip(req)
            rtrc <- rtr
        }()

        select {
        case rtr := <-rtrc:
            res = rtr.res
            err = rtr.err
        case <-c:
            err = errors.New("Exited before response could be obtained")
            return
        }

        if err == nil {
            return
        }

        select {
        case <-time.After(delay):
            n++
        case <-c:
            return
        }
    }

    return
}
