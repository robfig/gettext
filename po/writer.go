package po

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

// writer formats message fields into a buffer and writes to a destination.
// it is a mirror of the scanner.
type writer struct {
	buf *bytes.Buffer
	n   int64
}

func newWriter() writer {
	return writer{new(bytes.Buffer), 0}
}

// mul writes the given values on multiple lines, one per line.
func (wr *writer) mul(prefix string, vals []string) {
	for _, val := range vals {
		wr.buf.WriteString(prefix + val + "\n")
	}
}

// spc writes the given values on a single line, separated by spaces.
func (wr *writer) spc(prefix string, vals []string) {
	if len(vals) == 0 {
		return
	}
	wr.buf.WriteString(prefix)
	for i, val := range vals {
		if i > 0 {
			wr.buf.WriteString(" ")
		}
		wr.buf.WriteString(val)
	}
	wr.buf.WriteString("\n")
}

// one writes the given value with the given prefix.
func (wr *writer) one(prefix, val string) {
	if val != "" {
		wr.buf.WriteString(prefix + val + "\n")
	}
}

// opt writes the given value as a quoted string
func (wr *writer) opt(prefix, val string) {
	if val != "" {
		wr.quo(prefix, val)
	}
}

// quo always writes the given value (quoted), even if empty.
// Additionally, it breaks multiline strings across lines.
func (wr *writer) quo(prefix, val string) {
	if !strings.Contains(val, "\n") {
		wr.buf.WriteString(prefix + strconv.Quote(val) + "\n")
		return
	}

	// multiline
	wr.buf.WriteString(prefix + `""` + "\n")
	for {
		i := strings.Index(val, "\n")
		if i == -1 {
			if val != "" {
				wr.buf.WriteString(strconv.Quote(val) + "\n")
			}
			return
		}
		wr.buf.WriteString(strconv.Quote(val[:i+1]) + "\n")
		val = val[i+1:]
	}
}

// msgstr writes a singular msgstr.
// vals should be at most length 1.
func (wr *writer) msgstr(vals []string) {
	if len(vals) == 0 {
		wr.quo("msgstr ", "")
	} else {
		wr.quo("msgstr ", vals[0])
	}
}

// plural writes the plural form of msgstr.
func (wr *writer) plural(vals []string) {
	if len(vals) == 0 {
		wr.quo("msgstr[0] ", "")
	} else {
		for i, str := range vals {
			wr.quo("msgstr["+strconv.Itoa(i)+"] ", str)
		}
	}
}

// newline writes a newline
func (wr *writer) newline() {
	wr.buf.WriteString("\n")
}

// from writes to this buffer from the given WriterTo.
// it ignores errors, as writes to bytes.Buffer do not return errors.
func (wr *writer) from(w io.WriterTo) {
	w.WriteTo(wr.buf)
}

// to writes the contents of the writer to the given output.
func (wr *writer) to(w io.Writer) (n int64, err error) {
	return io.Copy(w, wr.buf)
}
