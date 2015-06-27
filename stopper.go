package stopper

import (
	"os"
	"os/signal"
)

type Stopper struct {
	Shutdown chan struct{}
	FatalErrors chan error
	Interrupts chan os.Signal
	Error error
}

func Create() *Stopper {
	
	stopper := Stopper{
		Shutdown	: make(chan struct{}),
		FatalErrors : make(chan error),
		Interrupts	: make(chan os.Signal, 1),
		Error 		: nil,
	}

	signal.Notify(stopper.Interrupts, os.Interrupt)

	exiting := false
	kill := func() {
		if !exiting {
			exiting = true
			close(stopper.Shutdown)
		}
	}

	go func() {
		for {
			select {
			case <- stopper.Interrupts:
				kill()
			case err := <- stopper.FatalErrors:
				if stopper.Error == nil {
					stopper.Error = err	
				}
				kill()
			}
		}
	}()

	return &stopper

}