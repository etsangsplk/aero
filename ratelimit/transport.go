package ratelimit

import "net/http"

type KeyFunc func(*http.Request) string

// Host based rate limiting
type HttpTransport struct {
    Transport http.RoundTripper
    KeyFunc   KeyFunc

    p   RateLimiter
    q   map[string]*rateQueue
}

func NewHttpTransport(p RateLimiter) *HttpTransport {
    return &HttpTransport{
        Transport: http.DefaultTransport,
        p:         p,
        q:         make(map[string]*rateQueue, 10),
        KeyFunc: func(req *http.Request) string {
            return req.URL.Host
        },
    }
}

func (t *HttpTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
    q, ok := t.q[t.KeyFunc(req)]
    if !ok {
        q = &rateQueue{
            l:         t.p.Clone(),
            treq:      make(chan *treq, 100),
            Transport: t.Transport,
        }
        q.loop()
        t.q[req.URL.Host] = q
    }
    r := &treq{
        req:  req,
        tres: make(chan *tres),
    }
    q.treq <- r
    s := <-r.tres
    return s.res, s.err
}

type treq struct {
    req  *http.Request
    tres chan *tres
}

type tres struct {
    res *http.Response
    err error
}

type rateQueue struct {
    l         RateLimiter
    treq      chan *treq
    Transport http.RoundTripper
}

func (q *rateQueue) loop() {
    go func() {
        defer close(q.treq)
        for {
            select {
            case tr := <-q.treq:
                select {
                case <-q.l.Tick():
                    go func(tr *treq) {
                        ts := &tres{}
                        ts.res, ts.err = q.Transport.RoundTrip(tr.req)
                        tr.tres <- ts
                    }(tr)
                }
            }
        }
    }()
}
