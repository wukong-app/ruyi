package engine

import (
	"math"
	"sync"
	"sync/atomic"

	"github.com/wukong-app/ruyi/pkg/base/contract"
	"github.com/wukong-app/ruyi/pkg/base/exception"
)

var _ contract.Ruyi = (*Ruyi)(nil)

func NewRuyi() *Ruyi {
	return &Ruyi{
		mx:          sync.Mutex{},
		description: "The Ruyi Jingu Bang, Sun Wukong’s magic staff, weighs thirteen thousand five hundred jin.",
		size:        20,
	}
}

// Ruyi is the implementation of ruyi.Ruyi
type Ruyi struct {
	mx          sync.Mutex
	description string
	size        int32
}

func (s *Ruyi) GetDescription() string {
	return s.description
}

func (s *Ruyi) GetSize() int32 {
	return s.size
}

func (s *Ruyi) Expand() (int32, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.size >= math.MaxInt32 {
		return 0, exception.ErrRuyiIsBigEnough
	}

	return atomic.AddInt32(&s.size, 1), nil
}

func (s *Ruyi) Shrink() (int32, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.size == 0 {
		return 0, exception.ErrRuyiIsSmallEnough
	}

	return atomic.AddInt32(&s.size, -1), nil
}
