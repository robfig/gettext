package po

import (
	"bytes"
	"io"
	"strconv"
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

// quo writes the given value as a quoted string
// TODO: break long strings across lines
func (wr *writer) quo(prefix, val string) {
	if val != "" {
		wr.str(prefix, val)
	}
}

// str always writes the given value (quoted), even if strty.
func (wr *writer) str(prefix, val string) {
	wr.buf.WriteString(prefix + strconv.Quote(val) + "\n")
}

// msgstr writes a singular msgstr.
// vals should be at most length 1.
func (wr *writer) msgstr(vals []string) {
	if len(vals) == 0 {
		wr.str("msgstr ", "")
	} else {
		wr.str("msgstr ", vals[0])
	}
}

// plural writes the plural form of msgstr.
func (wr *writer) plural(vals []string) {
	if len(vals) == 0 {
		wr.str("msgstr[0] ", "")
	} else {
		for i, str := range vals {
			wr.str("msgstr["+strconv.Itoa(i)+"] ", str)
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
