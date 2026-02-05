//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/wukong-app/ruyi/pkg/contract"
)

// New returns a new Ruyi.
func New() (contract.Ruyi, error) {
	panic(wire.Build(
		providerSet,
	))
}
