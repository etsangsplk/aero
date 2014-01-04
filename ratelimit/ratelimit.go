package ratelimit

import (
    "time"
)

type RateLimiter interface {
    Tick() <-chan time.Time
    Stop()
}

type Bursty struct {
    c   chan time.Time
    n   int
    t   time.Duration
    cls chan bool
}

func NewBursty(num int, unit time.Duration, burst int) *Bursty {
    l := &Bursty{
        c:   make(chan time.Time, burst),
        cls: make(chan bool, 1),
        n:   num,
        t:   unit,
    }

    l.start(burst)

    return l
}

func (l *Bursty) Tick() <-chan time.Time {
    return l.c
}

func (l *Bursty) Stop() {
    l.cls <- true
}

func (l *Bursty) start(burst int) {
    for i := 0; i < burst; i++ {
        l.c <- time.Now()
    }

    go func() {
        delta := l.t / time.Duration(l.n)
        ticker := time.NewTicker(delta)
        defer ticker.Stop()

        for {
            select {
            case now := <-ticker.C:
                select {
                case l.c <- now:
                default:
                }
            case <-l.cls:
                return
            }
        }
    }()
}
