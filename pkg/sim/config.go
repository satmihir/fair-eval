package sim

import "github.com/satmihir/fair/pkg/config"

// The global config comprising the simulation parameters including
// FAIR library turning params and traffic definition etc.
type SimulationConfig struct {
	// FAIR library tuning parameters
	FairConfig *config.FairnessTrackerConfig

	// Additinal FAIR parameters

	// Simulation parameters
	RegenerationRatePerSecond               float64   // How often the resource under contention is regenerated
	ClientRequestInterarrivalTimeMinSeconds []float64 // The interarrival time between requests (floor)
	ClientRequestInterarrivalTimeMaxSeconds []float64 // The interarrival time between requests (ceiling)
	ClientRequestStartTimeUnixMillis        []uint64  // The time when the clients start sending their requests
	ClientRequestEndTimeUnixMillis          []uint64  // The time when the clients stop sending their requests

	// Logs
	LogLocation string
}
