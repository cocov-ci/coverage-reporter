package test_helpers

import (
	"flag"
	"github.com/urfave/cli/v2"
	"os"
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

func GitHead(t *testing.T) string {
	c := exec.Command("git", "rev-parse", "HEAD")
	data, err := c.CombinedOutput()
	require.NoError(t, err)
	return strings.TrimSpace(string(data))
}

func MakeContext(t *testing.T, flags map[string]string) *cli.Context {
	fs := flag.FlagSet{}
	fs.String("token", "", "")

	require.NoError(t, fs.Set("token", "token"))

	for k, v := range flags {
		fs.String(k, "", "")
		require.NoError(t, fs.Set(k, v))
	}

	c := cli.NewContext(nil, &fs, nil)
	return c
}

func InDir(t *testing.T, path string, fn func(t *testing.T)) {
	pwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(path))
	defer func(pwd string) {
		require.NoError(t, os.Chdir(pwd))
	}(pwd)

	fn(t)
}

func MkTempDir(t *testing.T) string {
	d, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	return d
}
