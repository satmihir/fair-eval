package experiments

import (
	"context"
	"fmt"

	"github.com/satmihir/fair-eval/pkg/sim"
	"github.com/satmihir/fair/pkg/config"
)

const (
	HR      = 60*60*10 ^ 3
	TOT_LEN = 60*60*60*10 ^ 3
)

func RunGradualLoadExperiment() {
	conf := config.DefaultFairnessTrackerConfig()
	conf.Pd = .00005
	conf.IncludeStats = true
	fmt.Printf("conf.L: %d\n", conf.L)

	simConf := &sim.SimulationConfig{
		FairConfig:                conf,
		RegenerationRatePerSecond: 10,
		ClientRequestStartTimeUnixMillis: []uint64{
			0 * HR, 8 * HR, 12 * HR, 18 * HR, 18 * HR,
		},
		ClientRequestEndTimeUnixMillis: []uint64{
			TOT_LEN, 52 * HR, 48 * HR, 42 * HR, 42 * HR,
		},
		ClientRequestInterarrivalTimeMinSeconds: []float64{
			.1, .2, .5, .5, .25,
		},
		ClientRequestInterarrivalTimeMaxSeconds: []float64{
			.1, .2, .5, .5, .25,
		},
		LogLocation: "/tmp/gradual_load.log",
	}

	tr := sim.NewTrafficGen(simConf)
	expr := sim.NewExpr(simConf, tr)

	expr.Run(context.TODO())
}
