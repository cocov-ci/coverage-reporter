package test_helpers

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Timeout(t *testing.T, d time.Duration, fn func()) {
	ok := make(chan bool)
	go func() {
		fn()
		close(ok)
	}()

	select {
	case <-ok:
		return
	case <-time.After(d):
		t.Error("Run time exceeded")
	}
}

func GitRoot(t *testing.T) string {
	c := exec.Command("git", "rev-parse", "--show-toplevel")
	data, err := c.CombinedOutput()
	require.NoError(t, err)
	return strings.TrimSpace(string(data))
}
