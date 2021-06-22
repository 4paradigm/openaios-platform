package server

import (
	"github.com/4paradigm/openaios-platform/src/webterminal/gotty/webtty"
)

// Slave is webtty.Slave with some additional methods.
type Slave interface {
	webtty.Slave

	Close() error
}

type Factory interface {
	Name() string
	New(params map[string][]string) (Slave, error)
}
