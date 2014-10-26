package po

import "io"

// File represents a PO file.
type File struct {
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
	return File{msgs}, scan.Err()
}

// Write the PO file to a destination writer.
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	var wr = newWriter()
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
	wr.quo("msgctxt ", m.Ctxt)
	wr.quo("msgid ", m.Id)
	wr.quo("msgid_plural ", m.IdPlural)
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
