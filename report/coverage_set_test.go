package report

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoverageSet(t *testing.T) {
	set := CoverageSet{}
	set.SetLine(1, CoverageKindCovered, 5)
	set.SetLine(1, CoverageKindCovered, 5)
	set.SetLine(2, CoverageKindCovered, 20)
	set.SetLine(4, CoverageKindIgnore, 0)
	set.SetLine(5, CoverageKindCovered, 0)

	// Expected output is 10 (ascii) | separator | 20 (ascii) | neutral | ignore | missed

	set.FillMissingLines(CoverageKindNeutral)

	data := set.Encode()
	assert.Equal(t, "MTAeMjAAGxU=", data)
}

func TestCoverageSetLineZero(t *testing.T) {
	set := CoverageSet{}
	assert.PanicsWithError(t, "invalid zero line", func() {
		set.SetLine(0, CoverageKindCovered, 1)
	})
}

func TestCoverageSetFillMissingAsCovered(t *testing.T) {
	set := CoverageSet{}
	assert.PanicsWithError(t, "cannot fill missing lines as covered", func() {
		set.FillMissingLines(CoverageKindCovered)
	})
}
