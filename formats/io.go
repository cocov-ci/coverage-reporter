package formats

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func readFilePartialBuf(path string, offsetStart int, into []byte) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	defer func(f *os.File) { _ = f.Close() }(f)

	buf := bufio.NewReader(f)
	if offsetStart > 0 {
		discarded := 0
		for discarded < offsetStart {
			d, err := buf.Discard(offsetStart - discarded)
			if err != nil {
				return 0, err
			}
			discarded += d
		}
	}

	return io.ReadFull(buf, into)
}

func readFilePartial(path string, offsetStart, length int) ([]byte, error) {
	buf := make([]byte, length)
	i, err := readFilePartialBuf(path, offsetStart, buf)
	if err != nil {
		return nil, err
	}
	return buf[:i], nil
}

func bufferedLineReader(r io.Reader) *lineReader {
	return &lineReader{
		tmp: make([]byte, 0, 1024),
		r:   r,
	}
}

type lineReader struct {
	tmp []byte
	r   io.Reader
	eof bool
}

func (l *lineReader) lineFromTmp() (bool, string) {
	for {
		if len(l.tmp) == 0 {
			break
		}
		idx := bytes.IndexByte(l.tmp, '\n')
		if idx > -1 {
			str := strings.TrimSpace(string(l.tmp[:idx]))
			l.tmp = l.tmp[idx+1:]
			if len(str) == 0 {
				continue
			}
			return true, str
		}

		break
	}
	return false, ""
}

func (l *lineReader) NextLine() (string, error) {
	if ok, line := l.lineFromTmp(); ok {
		return line, nil
	}

	if l.eof {
		if len(l.tmp) == 0 {
			return "", io.EOF
		}

		str := string(l.tmp)
		l.tmp = l.tmp[:0]
		return strings.TrimSpace(str), nil
	}

	buf := make([]byte, 64)
	for {
		read, err := l.r.Read(buf)
		if err == io.EOF {
			if l.eof {
				str := strings.TrimSpace(string(l.tmp))
				l.tmp = l.tmp[:0]
				if len(str) == 0 {
					return "", io.EOF
				} else {
					return str, nil
				}
			}
			l.eof = true
		} else if err != nil {
			return "", err
		}

		l.tmp = append(l.tmp, buf[:read]...)
		if ok, line := l.lineFromTmp(); ok {
			return line, nil
		}
	}
}

func findFilePathPrefix(from, to string) *string {
	components := strings.Split(from, "/")
	var prefix []string
	if components[0] == "" {
		components = components[1:]
		prefix = append(prefix, "")
	}

	for len(components) > 0 {
		path := filepath.Join(to, strings.Join(components, "/"))
		s, err := os.Stat(path)
		if os.IsNotExist(err) {
			prefix = append(prefix, components[0])
			components = components[1:]
			continue
		} else if err != nil {
			continue
		}
		if s.IsDir() {
			prefix = append(prefix, components[0])
			components = components[1:]
			continue
		}

		prfx := strings.Join(prefix, "/") + "/"
		return &prfx
	}
	return nil
}
