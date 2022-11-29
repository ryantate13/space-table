package version

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("returns version from VERSION file", func(t *testing.T) {
		v, err := os.ReadFile("VERSION")
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(string(v)), Get())
	})
}
