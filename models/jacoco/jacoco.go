package jacoco

type Coverage struct {
	SessionInfo SessionInfo `xml:"sessioninfo"`
	Packages    []Package   `xml:"package"`
}

type SessionInfo struct {
	ID    string `xml:"id,attr"`
	Start string `xml:"start,attr"`
	Dump  string `xml:"dump,attr"`
}

type Package struct {
	Name        string       `xml:"name,attr"`
	SourceFiles []SourceFile `xml:"sourcefile"`
}

type SourceFile struct {
	Name  string `xml:"name,attr"`
	Lines []Line `xml:"line"`
}

type Line struct {
	Number int `xml:"nr,attr"`
	Hits   int `xml:"ci,attr"`
}
