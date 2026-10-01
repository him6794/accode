package executor

import "testing"

func TestLegacyExecutorIsDeprecated(t *testing.T) {
	t.Skip("legacy Go executor is deprecated; runtime is delegated to original Python back.py via handlers")
}
