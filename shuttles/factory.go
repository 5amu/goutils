package shuttles

import (
	"context"
	"fmt"
	"sync"
)

type ShuttleFactory struct {
	template       []string
	supply         []chan string
	lanes          int
	stopFactory    chan struct{}
	shuttleCounter int
	outputs        []*ShuttleOutput
}

func NewShuttleFactory(template []string, lanes int) *ShuttleFactory {
	var supply []chan string
	for i := 0; i < lanes-1; i++ {
		supply = append(supply, make(chan string, 1))
	}
	return &ShuttleFactory{
		supply:      supply,
		lanes:       lanes,
		template:    template,
		stopFactory: make(chan struct{}, 1),
	}
}

func (f *ShuttleFactory) SupplyLane(s string, l int) error {
	if l < 0 || l > f.lanes {
		return fmt.Errorf("lane %d does not exist in the shuttle factory", l)
	}
	f.supply[l] <- s
	return nil
}

func (f *ShuttleFactory) Stop() {
	f.stopFactory <- struct{}{}
}

func (f *ShuttleFactory) Start(ctx context.Context, parallelism int) error {
	// Create a new context
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	guard := make(chan struct{}, parallelism)
	var wg sync.WaitGroup

	for f.shuttleCounter = 0; true; f.shuttleCounter++ {
		shuttle := NewShuttle(childCtx, f.template, f.shuttleCounter)

		var toinject []string
		for j := 0; j < len(f.supply); j++ {
			select {
			case injarg := <-f.supply[j]:
				toinject = append(toinject, injarg)
			case <-f.stopFactory:
				childCtx.Done()
				wg.Wait()
				return nil
			}
		}
		shuttle.InjectArguments(toinject)

		guard <- struct{}{}
		wg.Add(1)
		go func() {
			if out, err := shuttle.Launch(); err != nil {
				fmt.Printf("[shuttle %v] ERROR: %v", shuttle.id, err)
			} else {
				f.outputs = append(f.outputs, out)
			}
			wg.Done()
			<-guard
		}()
	}
	return nil
}

func (f *ShuttleFactory) GetShuttleOutputs() []*ShuttleOutput {
	return f.outputs
}

func (f *ShuttleFactory) GetShuttleOutput(shuttleID int) *ShuttleOutput {
	if shuttleID > f.shuttleCounter || shuttleID < 0 {
		return nil
	}
	return f.outputs[shuttleID]
}
