package report

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

type CoverageSet struct {
	lines []*CoverageLine
}

func (c *CoverageSet) ensureSize(line int) {
	if currentLen := len(c.lines); currentLen <= line+1 {
		newLines := make([]*CoverageLine, line+1)
		copy(newLines, c.lines)
		c.lines = newLines
	}
}

func (c *CoverageSet) SetLine(line int, kind CoverageKind, hits int) {
	if line == 0 {
		panic(fmt.Errorf("invalid zero line"))
	}

	line = line - 1
	c.ensureSize(line)
	if c.lines[line] == nil {
		c.lines[line] = &CoverageLine{
			Kind: kind,
			Hits: hits,
		}
	} else {
		c.lines[line].Hits += hits
	}
}

func (c *CoverageSet) FillMissingLines(kind CoverageKind) {
	if kind == CoverageKindCovered {
		panic(fmt.Errorf("cannot fill missing lines as covered"))
	}

	for i, v := range c.lines {
		if v == nil {
			c.lines[i] = &CoverageLine{Kind: kind}
		}
	}
}

const (
	separator byte = 0x1E
	neutral   byte = 0x00
	ignore    byte = 0x1B
	missed    byte = 0x15
)

func (c *CoverageSet) Encode() string {
	buf := bytes.Buffer{}
	lastWasAscii := false
	for _, v := range c.lines {
		switch v.Kind {
		case CoverageKindCovered:
			if v.Hits > 0 {
				if lastWasAscii {
					buf.WriteByte(separator)
				}
				buf.Write([]byte(fmt.Sprintf("%d", v.Hits)))
				lastWasAscii = true
			} else {
				buf.WriteByte(missed)
				lastWasAscii = false
			}
		case CoverageKindIgnore:
			buf.WriteByte(ignore)
			lastWasAscii = false
		case CoverageKindNeutral:
			buf.WriteByte(neutral)
			lastWasAscii = false
		}
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
