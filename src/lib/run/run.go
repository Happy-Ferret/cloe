package run

import (
	"fmt"
	"os"
	"sync"

	"github.com/tisp-lang/tisp/src/lib/compile"
	"github.com/tisp-lang/tisp/src/lib/core"
)

const maxConcurrentOutputs = 256

var sem = make(chan bool, maxConcurrentOutputs)

// Run runs outputs.
func Run(os []compile.Output) {
	wg := sync.WaitGroup{}

	for _, v := range os {
		wg.Add(1)

		if v.Expanded() {
			go func() {
				evalOutputList(v.Value())
				wg.Done()
			}()
		} else {
			sem <- true
			go runOutput(v.Value(), &wg)
		}
	}

	wg.Wait()
}

func evalOutputList(t *core.Thunk) {
	wg := sync.WaitGroup{}

	for {
		v := core.PApp(core.Equal, t, core.EmptyList).Eval()

		if b, ok := v.(core.BoolType); !ok {
			failOnError(core.NotBoolError(v).Eval().(core.ErrorType))
		} else if b {
			break
		}

		wg.Add(1)
		sem <- true
		go runOutput(core.PApp(core.First, t), &wg)

		t = core.PApp(core.Rest, t)
	}

	wg.Wait()
}

func runOutput(t *core.Thunk, wg *sync.WaitGroup) {
	if err, ok := t.EvalOutput().(core.ErrorType); ok {
		failOnError(err)
	}

	<-sem
	wg.Done()
}

func failOnError(err core.ErrorType) {
	fmt.Fprint(os.Stderr, err.Lines())
	os.Exit(1)
}
