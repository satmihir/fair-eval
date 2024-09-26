package sim

import (
	"context"
	"fmt"
	"log"
	"os"

	stime "github.com/satmihir/fair-eval/pkg/time"
	"github.com/satmihir/fair/pkg/request"
	"github.com/satmihir/fair/pkg/tracker"
)

var ErrNoTokens = fmt.Errorf("no tokens")

type TokenBucket struct {
	tokens                float64
	tokensPerSecond       float64
	lastUpdatedTimeMillis uint64
	clock                 *stime.SimClock
}

func NewTokenBucket(initialTokens uint32, tokensPerSecond float64, clock *stime.SimClock) *TokenBucket {
	return &TokenBucket{
		tokens:                float64(initialTokens),
		tokensPerSecond:       tokensPerSecond,
		lastUpdatedTimeMillis: 0,
		clock:                 clock,
	}
}

func (tb *TokenBucket) Take() error {
	currentTimeMillis := uint64(tb.clock.Now().UnixMilli())
	diff := float64(currentTimeMillis-tb.lastUpdatedTimeMillis) / 1000.0

	tb.tokens += tb.tokensPerSecond * float64(diff)
	tb.lastUpdatedTimeMillis = currentTimeMillis

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return nil
	}

	return ErrNoTokens
}

type Expr struct {
	config    *SimulationConfig
	traffic   *TrafficGen
	tracker   *tracker.FairnessTracker
	bkt       *TokenBucket
	simClock  *stime.SimClock
	simTicker *stime.SimTicker
	logFile   *os.File
}

func NewExpr(config *SimulationConfig, traffic *TrafficGen) *Expr {
	clk := stime.NewSimClock()
	tk := stime.NewNeverTicker()
	trk, err := tracker.NewFairnessTrackerWithClockAndTicker(config.FairConfig, clk, tk)
	if err != nil {
		log.Fatalf("failure when building tracker %v", err)
	}

	file, err := os.OpenFile(config.LogLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	failOnError(err)

	return &Expr{
		config:    config,
		traffic:   traffic,
		tracker:   trk,
		bkt:       NewTokenBucket(uint32(config.RegenerationRatePerSecond), config.RegenerationRatePerSecond, clk),
		simClock:  clk,
		simTicker: tk,
		logFile:   file,
	}
}

func (e *Expr) Run(ctx context.Context) {
	//prevTicker := e.simClock.Now()
	for {
		nxt, err := e.traffic.NextRequest()
		if err == ErrEoq {
			return
		} else if err != nil {
			failOnError(err)
		}

		e.simClock.SetTimeMillis(nxt.timeMillis)
		//if e.simClock.Now().Sub(prevTicker) >= 5*time.Minute {
		//	e.simTicker.SendTick(e.simClock.Now())
		//}

		//prevTicker = e.simClock.Now()

		resp, err := e.tracker.RegisterRequest(ctx, []byte(nxt.clientId))
		failOnError(err)

		if resp.ShouldThrottle {
			// The request got throttled
			e.logFile.WriteString(fmt.Sprintf("%s,%d,T\n", nxt.clientId, nxt.timeMillis))
			fmt.Printf("T %s [%v]\n", nxt.clientId, resp.ResultStats.BucketProbabilities)
			continue
		}

		if err := e.bkt.Take(); err == nil {
			e.logFile.WriteString(fmt.Sprintf("%s,%d,S\n", nxt.clientId, nxt.timeMillis))
			fmt.Printf("S %s [%v]\n", nxt.clientId, resp.ResultStats.BucketProbabilities)
			e.tracker.ReportOutcome(ctx, []byte(nxt.clientId), request.OutcomeSuccess)
		} else {
			e.logFile.WriteString(fmt.Sprintf("%s,%d,F\n", nxt.clientId, nxt.timeMillis))
			fmt.Printf("F %s [%v]\n", nxt.clientId, resp.ResultStats.BucketProbabilities)
			e.tracker.ReportOutcome(ctx, []byte(nxt.clientId), request.OutcomeFailure)
		}
	}
}

func failOnError(err error) {
	if err != nil {
		log.Fatalf("Unexpected error %v", err)
	}
}
