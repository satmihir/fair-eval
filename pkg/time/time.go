package time

import (
	"log"
	"time"
)

// A simulation clock
type SimClock struct {
	// Current (simulation) time in epoch millis. Starts at 0.
	curTimeMillis uint64
}

func NewSimClock() *SimClock {
	return &SimClock{
		curTimeMillis: 0,
	}
}

func (c *SimClock) Now() time.Time {
	return timeFromUnixMillis(int64(c.curTimeMillis))
}

func (c *SimClock) Sleep(duration time.Duration) {
	log.Fatalf("not implemented")
}

func (c *SimClock) SetTimeMillis(timeInMillis uint64) {
	c.curTimeMillis = timeInMillis
}

// A simulated ticker that never actually ticks
type SimTicker struct {
	channel chan time.Time
}

func NewNeverTicker() *SimTicker {
	return &SimTicker{
		channel: make(chan time.Time),
	}
}

func (t *SimTicker) SendTick(tm time.Time) {
	t.channel <- tm
}

func (t *SimTicker) C() <-chan time.Time {
	return t.channel
}

func (t *SimTicker) Stop() {
	// do nothing
}

// Create a time object from the given unix millis
func timeFromUnixMillis(unixMillis int64) time.Time {
	u := float64(unixMillis) / 1000.
	if unixMillis < 1000 {
		u = 0
	}
	return time.Unix(int64(u), (unixMillis%10^3)*10^6)
}
