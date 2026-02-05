package ruyi

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wukong-app/ruyi/internal/core"
)

// Test_NormalizeConcept test core.NormalizeConcept
func Test_NormalizeConcept(t *testing.T) {
	t.Run("not_exist", func(t *testing.T) {
		concept, exist := core.NormalizeConcept("this_is_not_exist_key")
		require.False(t, exist)
		require.Equal(t, core.Concept{}, concept)
	})
	t.Run("exist", func(t *testing.T) {
		concept, exist := core.NormalizeConcept("jpg")
		require.True(t, exist)
		require.Equal(t, core.JPEG(), concept)
	})
}
