package report

type CoverageKind int

const (
	CoverageKindNeutral CoverageKind = iota
	CoverageKindIgnore
	CoverageKindCovered
)

type CoverageLine struct {
	Kind CoverageKind
	Hits int
}
