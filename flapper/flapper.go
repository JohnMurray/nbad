package flapper

import (
	"time"

	"github.com/JohnMurray/nbad/timewindow"
)

// Flapper is a simple struct for watching for services that are in a flapState
type Flapper struct {
	// the max amount of state-transitions that can happen before "flapping"
	max uint

	// size of sliding time-window for state-change counters in seconds
	duration uint

	// a map of servie-names to sliding window-counters
	services map[string]*timewindow.Window
}

// NewFlapper - Create a new instance of Flapper
func NewFlapper(max uint, duration uint) *Flapper {
	f := &Flapper{
		max:      max,
		duration: duration,
		services: make(map[string]*timewindow.Window),
	}
	return f
}

// NoteStateChange -
// Increment the counter for a service (or create a counter for the service if one has not
// alredy been created). (Lazily create services).
func (f *Flapper) NoteStateChange(service string) {
	if state, ok := f.services[service]; ok {
		state.Add(time.Now().Unix(), 1)
	} else {
		f.services[service] = timewindow.New(time.Now().Unix(), int(f.duration))
		f.services[service].Add(time.Now().Unix(), 1)
	}
}

// IsFlapping -
// Return boolean indicating whether or not the state is flapping or not. Note that if the
// state does not exist we always return false sine we lazily create counters.
//
// service   - the service to check against
// recompute - bool flag to recompute the time-window if it has not been updated in a while
func (f *Flapper) IsFlapping(service string, recompute bool) bool {
	if state, ok := f.services[service]; ok {
		if recompute {
			state.Add(time.Now().Unix(), 0)
		}
		return state.Total() >= int(f.max)
	}
	return false
}

// Compact -
// If a Flapper has been running for a long time, you may want to periodically clean up any
// services that do not have any data. Compaction allows us to compress our internal data-structures
// and potentially free up memory.
func (f *Flapper) Compact() {
	for service, counter := range f.services {
		if counter.Total() == 0 {
			delete(f.services, service)
		}
	}
}
