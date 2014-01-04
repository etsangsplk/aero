package backoff

import (
    "math/rand"
    "time"
)

type Backoffer interface {
    // The next delay and whether or not to try agian
    //
    // When the FailAfter retries is reached, try again is false and delay is 0
    Next(r uint) (delay time.Duration, tryAgain bool)

    // reset the interal counter
    Reset()

    // number of retries
    //
    // 0 for unbounded retries, not recommended
    // Default is 20 retries
    FailAfter(retries uint)

    // Fractional amount (jitter / delay) to randomly vary the time before and after the mathematical delay
    //
    // For example, a with Jitter(0.1, 0.2) if the next delay is 10 min (calculated by the mathematical rule) then the actual delay could be between 9 min and 12 min
    Jitter(before float64, after float64)

    // The maximum delay
    // 0 for unbounded delay
    //
    // Default is 5 minutes
    MaxDelay(d time.Duration)

    // The delay for the first Next() call
    //
    // Default is 100 ms
    InitialDelay(d time.Duration)
}

type ExponentialBackoff struct {
    failAfter      uint
    jitter, offset float64
    r              uint
    rand           func() float64

    initialDelay, maxDelay time.Duration
}

func NewExponential() *ExponentialBackoff {
    rand.Seed(time.Now().UTC().UnixNano())
    return &ExponentialBackoff{
        failAfter:    20,
        jitter:       0,
        offset:       0,
        maxDelay:     5 * time.Minute,
        initialDelay: 100 * time.Millisecond,
        rand:         rand.Float64,
        r:            0,
    }
}

func (b *ExponentialBackoff) Next(r uint) (delay time.Duration, tryAgain bool) {
    tryAgain = true

    if r == 0 {
        delay = b.initialDelay
    } else if r < b.failAfter {
        delay = b.initialDelay * (2 << (r - 1))
    } else {
        delay = 0
        tryAgain = false
    }

    if delay > b.maxDelay {
        delay = b.maxDelay
    }

    if tryAgain && b.jitter != 0.0 {
        j := 1.0 + b.rand()*b.jitter + b.offset
        delay = time.Duration(float64(delay) * j)
    }

    return
}

func (b *ExponentialBackoff) Reset() {
    b.r = 0
}

func (b *ExponentialBackoff) FailAfter(retries uint) {
    b.failAfter = retries
}

func (b *ExponentialBackoff) MaxDelay(d time.Duration) {
    b.maxDelay = d
}

func (b *ExponentialBackoff) InitialDelay(d time.Duration) {
    b.initialDelay = d
}
func (b *ExponentialBackoff) Jitter(before float64, after float64) {
    b.jitter = before + after
    b.offset = -1.0 * before
}
