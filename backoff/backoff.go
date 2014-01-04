package backoff

import (
    "math/rand"
    "time"
)

type Backoffer interface {
    // The next delay and whether or not to try agian
    //
    // When the FailAfter retries is reached, try again is false and delay is 0
    Next(fails uint) (delay time.Duration, tryAgain bool)
}

type Exp struct {
    // Number of retries
    //
    // 0 for unbounded retries, not recommended
    // Default is 10 retries
    FailAfter uint

    // Fractional amount (of delay) to randomly vary the time before and after the mathematical delay
    //
    // For example, a with Jitter(0.1, 0.2) if the next delay is 10 min (calculated by the mathematical rule) then the actual delay could be between 9 min and 12 min
    JitterBefore, JitterAfter float64

    // The maximum delay
    //
    // 0 for unbounded delay
    // Default is 1 minutes
    MaxDelay time.Duration

    // The delay for the first Next() call
    //
    // Default is 100 ms
    InitialDelay time.Duration

    rand func() float64
}

func NewExp() *Exp {
    rand.Seed(time.Now().UTC().UnixNano())
    return &Exp{
        FailAfter:    10,
        MaxDelay:     time.Minute,
        InitialDelay: 100 * time.Millisecond,
        rand:         rand.Float64,
    }
}

func (b *Exp) Next(fails uint) (delay time.Duration, tryAgain bool) {
    tryAgain = true

    if fails == 0 {
        delay = b.InitialDelay
    } else if fails < b.FailAfter {
        delay = b.InitialDelay * (2 << (fails - 1))
    } else {
        delay = 0
        tryAgain = false
    }

    if fails == b.FailAfter-1 {
        delay = 0
    }

    if delay > b.MaxDelay {
        delay = b.MaxDelay
    }

    if tryAgain && (b.JitterBefore+b.JitterAfter) != 0.0 {
        j := 1.0 + b.rand()*(b.JitterAfter+b.JitterBefore) - b.JitterBefore
        delay = time.Duration(float64(delay) * j)
    }

    return
}
