package po

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

// scanner scans the message fields of a po file.
// it is a mirror of the writer.
type scanner struct {
	*bufio.Scanner
	hasNext bool
	err     error
}

func newScanner(r io.Reader) *scanner {
	s := &scanner{bufio.NewScanner(r), true, nil}
	s.next()
	return s
}

// next returns true if there is another record to read.
func (s *scanner) next() {
	// Skip blank lines between records.
	s.hasNext = s.Scan()
	for s.hasNext && len(bytes.TrimSpace(s.Bytes())) == 0 {
		s.hasNext = s.Scan()
	}
}

func (s *scanner) mul(prefix string) []string {
	var r []string
	for s.hasNext && bytes.HasPrefix(s.Bytes(), []byte(prefix)) {
		r = append(r, s.Text()[len(prefix):])
		s.next()
	}
	return r
}

func (s *scanner) spc(prefix string) []string {
	var r []string
	if s.hasNext && bytes.HasPrefix(s.Bytes(), []byte(prefix)) {
		r = append(r, strings.Fields(s.Text()[len(prefix):])...)
		s.next()
	}
	return r
}

func (s *scanner) one(prefix string) string {
	var r string
	if s.hasNext && bytes.HasPrefix(s.Bytes(), []byte(prefix)) {
		r = s.Text()[len(prefix):]
		s.next()
	}
	return r
}

func (s *scanner) quo(prefix string) string {
	var r string
	if s.hasNext && bytes.HasPrefix(s.Bytes(), []byte(prefix)) {
		r = s.unquote(s.Text()[len(prefix):])
		s.next()
	}
	return r
}

func (s *scanner) msgstr() []string {
	var r []string
	for s.hasNext && bytes.HasPrefix(s.Bytes(), []byte("msgstr")) {
		r = append(r, s.unquote(s.Text()[len("msgstr "):]))
		s.next()

		// TODO: Plural
		// var i int
		// var str string
		// _, err = fmt.Scanf("msgstr[%d] %q", &i, &str)
		// if err != nil {
		// }
	}
	return r
}

func (s *scanner) unquote(str string) string {
	var r, err = strconv.Unquote(str)
	if err != nil {
		s.err = err
		s.hasNext = false
	}
	return r
}

func (s *scanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.Scanner.Err()
}
