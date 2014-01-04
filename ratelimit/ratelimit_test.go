package ratelimit

import (
    "testing"
    "time"
)

func TestBurstyRateLimit(t *testing.T) {
    const count = 4
    delta := 5 * time.Millisecond
    burst := 3

    limiter := NewBursty(1, delta, burst)

    t0 := time.Now()
    for i := 0; i < count; i++ {
        <-limiter.Tick()
    }
    t1 := time.Now()
    dt := t1.Sub(t0)

    limiter.Stop()

    target := delta * time.Duration(count-burst)
    slop := target * 2 / 10
    if dt < target-slop || (!testing.Short() && dt > target+slop) {
        t.Fatalf("%d %s ticks with a burst of %d took %s, expected [%s,%s]", count, delta, burst, dt, target-slop, target+slop)
    }

    // Now test that the ticker stopped
    select {
    case <-limiter.Tick():
        t.Fatal("Ticker did not shut down")
    default:
        // ok
    }
}

func TestBurstyRateLimitBuildup(t *testing.T) {
    const count = 4
    delta := 5 * time.Millisecond
    burst := 3

    limiter := NewBursty(1, delta, burst)

    // wait some time to make sure the channel doesn't build up more than burst
    time.Sleep(delta * 3 * 101 / 100)

    t0 := time.Now()
    for i := 0; i < count; i++ {
        <-limiter.Tick()
    }
    t1 := time.Now()
    dt := t1.Sub(t0)

    limiter.Stop()

    target := delta * time.Duration(count-burst)
    slop := target * 2 / 10
    if dt < target-slop || (!testing.Short() && dt > target+slop) {
        t.Fatalf("%d %s ticks with a burst of %d took %s, expected [%s,%s]", count, delta, burst, dt, target-slop, target+slop)
    }
}
