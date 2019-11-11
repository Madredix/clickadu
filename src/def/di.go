package def

import (
	"fmt"
	"log"
	"sync"

	"github.com/sarulabs/di"
)

var (
	mu       sync.Mutex
	context  di.Container
	builders []buildFn
)

type (
	// public
	Builder    = di.Builder
	Context    = di.Container
	Definition = di.Def

	// private
	buildFn func(builder *di.Builder, params map[string]interface{}) error
)

// Register definition builder
func Register(fn buildFn) {
	mu.Lock()
	defer mu.Unlock()

	builders = append(builders, fn)
}

// Get context
func Instance(params map[string]interface{}) (di.Container, error) {
	if context != nil {
		return context, nil
	}

	builder, err := di.NewBuilder()
	if err != nil {
		return nil, fmt.Errorf("can't create context builder: %s", err)
	}

	for _, fn := range builders {
		if err := fn(builder, params); err != nil {
			return nil, err
		}
	}

	context = builder.Build()

	return context, nil
}

// Get context for tests
func TestInstance() di.Container {
	diContext, err := Instance(map[string]interface{}{
		"configFile": `./config/config_test.json`,
	})
	if err != nil {
		log.Fatalf(`error init di container: %s`, err)
	}

	return diContext
}
