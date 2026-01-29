package ruyi

import (
	"github.com/wukong-app/ruyi/internal/engine"
	"github.com/wukong-app/ruyi/pkg/base/contract"
)

// New returns a new Ruyi.
func New() contract.Ruyi {
	return engine.NewRuyi()
}
