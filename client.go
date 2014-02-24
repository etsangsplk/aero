package aero

import (
    "github.com/mceldeen/aero/backoff"
    "github.com/mceldeen/aero/ratelimit"
    "net"
    "net/http"
    "time"
)

func NewHTTPClient(backoffer backoff.Backoffer, limiter ratelimit.RateLimiter, timeout time.Duration) *http.Client {
    bt := backoff.NewHttpTransport(backoffer)
    bt.Transport = &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        ResponseHeaderTimeout: time.Duration(float64(timeout) * 1.5),
        Dial: func(network, addr string) (net.Conn, error) {
            return net.DialTimeout(network, addr, timeout)
        },
    }

    lt := ratelimit.NewHttpTransport(limiter)
    lt.Transport = bt

    return &http.Client{
        Transport: lt,
    }
}

func NewHTTPClientWithKeyFunc(backoffer backoff.Backoffer, limiter ratelimit.RateLimiter, timeout time.Duration, keyFunc ratelimit.KeyFunc) *http.Client {
    bt := backoff.NewHttpTransport(backoffer)
    bt.Transport = &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        ResponseHeaderTimeout: time.Duration(float64(timeout) * 1.5),
        Dial: func(network, addr string) (net.Conn, error) {
            return net.DialTimeout(network, addr, timeout)
        },
    }

    lt := ratelimit.NewHttpTransport(limiter)
    lt.Transport = bt
    lt.KeyFunc = keyFunc

    return &http.Client{
        Transport: lt,
    }
}
