package ruyi

import (
	"github.com/wukong-app/ruyi/internal"
	"github.com/wukong-app/ruyi/pkg/contract"
)

// New returns a new Ruyi.
func New() (contract.Ruyi, error) {
	return internal.New()
}
