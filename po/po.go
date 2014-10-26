package po

import (
	"bufio"
	"bytes"
	"io"
	"net/textproto"
	"sort"
	"strings"
)

// File represents a PO file.
type File struct {
	Header   textproto.MIMEHeader
	Messages []Message
}

// Message stores a gettext message.
type Message struct {
	Comment           // message comments from PO file
	Ctxt     string   // msgctxt: message context, if any
	Id       string   // msgid: untranslated singular string
	IdPlural string   // msgid_plural: untranslated plural string
	Str      []string // msgstr or msgstr[n]: translated strings
}

// Comment stores meta-data from a gettext message.
type Comment struct {
	TranslatorComments []string
	ExtractedComments  []string
	References         []string
	Flags              []string
	PrevCtxt           string
	PrevId             string
	PrevIdPlural       string
}

// Parse reads the content of a PO file and returns the list of messages.
func Parse(r io.Reader) (File, error) {
	var msgs []Message
	var scan = newScanner(r)
	for scan.nextmsg() {
		// NOTE: the source code order of these fields is important.
		var msg = Message{
			Comment: Comment{
				TranslatorComments: scan.mul("# "),
				ExtractedComments:  scan.mul("#."),
				References:         scan.spc("#:"),
				Flags:              scan.spc("#,"),
				PrevCtxt:           scan.one("#| msgctxt"),
				PrevId:             scan.one("#| msgid"),
				PrevIdPlural:       scan.one("#| msgid_plural"),
			},
			Ctxt:     scan.quo("msgctxt"),
			Id:       scan.quo("msgid"),
			IdPlural: scan.quo("msgid_plural"),
			Str:      []string{scan.quo("msgstr")},
		}
		msgs = append(msgs, msg)
	}
	if scan.Err() != nil {
		return File{}, scan.Err()
	}

	var header textproto.MIMEHeader
	if msgs[0].Id == "" && len(msgs[0].Str) == 1 {
		var err error
		header, err = textproto.NewReader(bufio.NewReader(strings.NewReader(msgs[0].Str[0]))).
			ReadMIMEHeader()
		if err != nil && err != io.EOF {
			return File{}, err
		}
		msgs = msgs[1:]
	}

	return File{header, msgs}, nil
}

// Write the PO file to a destination writer.
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	var wr = newWriter()
	// TODO: Probably better to make a type for the header and implement WriterTo
	if len(f.Header) > 0 {
		wr.quo("msgid ", "")
		var keys []string
		for k := range f.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var buf bytes.Buffer
		for _, k := range keys {
			buf.WriteString(k + ": " + f.Header.Get(k) + "\n")
		}
		wr.quo("msgstr ", buf.String())
		wr.newline()
	}
	for _, msg := range f.Messages {
		wr.from(msg)
		wr.newline()
	}
	return wr.to(w)
}

// Write the PO Message to a destination writer.
func (m Message) WriteTo(w io.Writer) (n int64, err error) {
	var wr = newWriter()
	wr.from(m.Comment)
	wr.opt("msgctxt ", m.Ctxt)
	wr.quo("msgid ", m.Id)
	wr.opt("msgid_plural ", m.IdPlural)
	wr.msgstr(m.Str)

	// TODO: If there is a plural form specified, then msgstr has an index.
	// if len(m.IdPlural) == 0 {
	// wr.msgstr(m.Str)
	// } else {
	// 	wr.plural(m.Str)
	// }

	return wr.to(w)
}

// Write the comment to the given writer.
func (c Comment) WriteTo(w io.Writer) (n int64, err error) {
	var wr = newWriter()
	wr.mul("#  ", c.TranslatorComments)
	wr.mul("#. ", c.ExtractedComments)
	wr.spc("#: ", c.References)
	wr.spc("#, ", c.Flags)
	wr.one("#| msgctxt ", c.PrevCtxt)
	wr.one("#| msgid ", c.PrevId)
	wr.one("#| msgid_plural ", c.PrevIdPlural)
	return wr.to(w)
}
