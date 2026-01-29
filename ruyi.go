package ruyi

import "github.com/wukong-app/ruyi/internal/engine"

// Ruyi is the interface for the Ruyi.
type Ruyi interface {

	// GetDescription returns the description of the Ruyi.
	// @return desc description
	GetDescription() (desc string)

	// GetSize returns the size of the Ruyi.
	// @return size current size
	GetSize() (size int32)

	// Expand returns the expanded size of the Ruyi.
	// @return size expanded size
	// @return err cannot continue expanding
	Expand() (size int32, err error)

	// Shrink returns the shrunk size of the Ruyi.
	// @return size shrunk size
	// @return err cannot continue shrinking
	Shrink() (size int32, err error)
}

// New returns a new Ruyi.
func New() Ruyi {
	return engine.NewRuyi()
}
