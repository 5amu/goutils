package shuttles

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type ShuttleFactory struct {
	template       []string
	supply         []chan string
	lanes          int
	shuttleCounter int
	stopFactory    chan struct{}
}

func NewShuttleFactory(template string, lanes int) *ShuttleFactory {
	var supply []chan string
	for i := 0; i < lanes; i++ {
		supply = append(supply, make(chan string, 1))
	}
	return &ShuttleFactory{supply: supply, lanes: lanes, shuttleCounter: 0, template: strings.Split(template, " ")}
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

	for i := 0; true; i++ {
		shuttle := NewShuttle(childCtx, f.template, i)

		var toinject []string
		for j := 0; j < len(f.supply); j++ {
			select {
			case injarg := <-f.supply[j]:
				toinject = append(toinject, injarg)
			case <-f.stopFactory:
				wg.Wait()
				return nil
			}
		}
		shuttle.InjectArguments(toinject)

		guard <- struct{}{}
		wg.Add(1)
		go func() {
			if err := shuttle.Launch(); err != nil {
				panic(err)
			}
			fmt.Println(shuttle.output.Output)
			wg.Done()
			<-guard
		}()
	}
	return nil
}
