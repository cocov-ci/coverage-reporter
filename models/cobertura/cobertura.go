package cobertura

type Coverage struct {
	Sources  []string  `xml:"sources>source"`
	Packages []Package `xml:"packages>package"`
}

func (c Coverage) Valid() bool {
	return len(c.Sources) > 0 && c.validatePackages()
}

func (c Coverage) validatePackages() bool {
	if len(c.Packages) == 0 {
		return true
	}

	for _, p := range c.Packages {
		if !p.valid() {
			return false
		}
	}

	return true
}

type Package struct {
	Name    string  `xml:"name,attr"`
	Classes []Class `xml:"classes>class"`
}

func (p Package) valid() bool {
	if len(p.Name) == 0 {
		return false
	}

	for _, c := range p.Classes {
		if !c.valid() {
			return false
		}
	}

	return true
}

type Class struct {
	Filename string `xml:"filename,attr"`
	Lines    []Line `xml:"lines>line"`
}

func (c Class) valid() bool {
	if len(c.Filename) == 0 {
		return false
	}

	for _, l := range c.Lines {
		if !l.valid() {
			return false
		}
	}

	return true
}

type Line struct {
	Number int `xml:"number,attr"`
	Hits   int `xml:"hits,attr"`
}

func (l Line) valid() bool {
	return l.Number > 0
}
