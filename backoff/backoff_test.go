package backoff

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestExpBackoff(t *testing.T) {
    delays := []time.Duration{
        10 * time.Millisecond,
        20 * time.Millisecond,
        40 * time.Millisecond,
        80 * time.Millisecond,
        160 * time.Millisecond,
        320 * time.Millisecond,
        640 * time.Millisecond,
        1000 * time.Millisecond,
        1000 * time.Millisecond,
        0 * time.Millisecond,
    }

    b := NewExp()
    b.InitialDelay = 10 * time.Millisecond
    b.MaxDelay = 1000 * time.Millisecond
    b.FailAfter = 10
    n := uint(0)
    for delay, next := b.Next(n); next; delay, next = b.Next(n) {
        assert.Equal(t, delays[n], delay, fmt.Sprintf("%d", n))
        n++
    }

    assert.Equal(t, uint(10), n)

    delayRanges := [][]time.Duration{
        []time.Duration{9 * time.Millisecond, 12 * time.Millisecond},
        []time.Duration{18 * time.Millisecond, 24 * time.Millisecond},
        []time.Duration{36 * time.Millisecond, 48 * time.Millisecond},
        []time.Duration{72 * time.Millisecond, 96 * time.Millisecond},
        []time.Duration{144 * time.Millisecond, 192 * time.Millisecond},
        []time.Duration{288 * time.Millisecond, 384 * time.Millisecond},
        []time.Duration{576 * time.Millisecond, 768 * time.Millisecond},
        []time.Duration{900 * time.Millisecond, 1200 * time.Millisecond},
        []time.Duration{900 * time.Millisecond, 1200 * time.Millisecond},
        []time.Duration{0 * time.Millisecond, 0 * time.Millisecond},
    }
    b.JitterBefore = 0.1
    b.JitterAfter = 0.2

    // min
    b.rand = func() float64 {
        return 0.0
    }
    n = 0
    for delay, next := b.Next(n); next; delay, next = b.Next(n) {
        assert.Equal(t, delayRanges[n][0], delay)
        n++
    }

    // max
    b.rand = func() float64 {
        return 1.0
    }
    n = 0
    for delay, next := b.Next(n); next; delay, next = b.Next(n) {
        assert.Equal(t, delayRanges[n][1], delay)
        n++
    }
}
