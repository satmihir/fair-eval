package sim

import (
	"fmt"
	"math/rand"
	"sort"
)

var ErrEoq = fmt.Errorf("requests ended")

type Request struct {
	clientId   string
	timeMillis uint64
}

func NewRequest(clientId string, timeMillis uint64) *Request {
	return &Request{
		clientId:   clientId,
		timeMillis: timeMillis,
	}
}

type TrafficGen struct {
	requestQueue []*Request
	queuePos     int
	config       *SimulationConfig
}

func NewTrafficGen(config *SimulationConfig) *TrafficGen {
	r := rand.New(rand.NewSource(42))
	q := make([]*Request, 0)

	// Generate and append all the client requests to the queue
	for cl := 0; cl < len(config.ClientRequestEndTimeUnixMillis); cl++ {
		curTime := config.ClientRequestStartTimeUnixMillis[cl]
		for curTime < config.ClientRequestEndTimeUnixMillis[cl] {
			min := config.ClientRequestInterarrivalTimeMinSeconds[cl] * 1000
			max := config.ClientRequestInterarrivalTimeMaxSeconds[cl] * 1000
			iatMillis := min + r.Float64()*(max-min)

			q = append(q, NewRequest(fmt.Sprintf("client-%d", cl), curTime+uint64(iatMillis)))
			curTime += uint64(iatMillis)
		}
	}

	// Sort all the events by time
	sort.Slice(q, func(i, j int) bool {
		return q[i].timeMillis < q[j].timeMillis
	})

	return &TrafficGen{
		requestQueue: q,
		queuePos:     0,

		config: config,
	}
}

func (tg *TrafficGen) NextRequest() (*Request, error) {
	defer func() {
		tg.queuePos++
	}()

	if tg.queuePos == len(tg.requestQueue) {
		return nil, ErrEoq
	}

	return tg.requestQueue[tg.queuePos], nil
}
