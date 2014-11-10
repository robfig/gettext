package po

import (
	"bytes"
	"net/textproto"
	"reflect"
	"strings"
	"testing"
)

var po = `
msgid ""
msgstr ""
"Content-Transfer-Encoding: 8bit\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Language: sk\n"
"Language-Team: Slovak <sk-i18n@lists.linux.sk>\n"
"Last-Translator: Marcel Telka <marcel@telka.sk>\n"
"Mime-Version: 1.0\n"
"Plural-Forms: nplurals=3; plural=(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2;\n"
"Po-Revision-Date: 2014-05-10 18:15+0200\n"
"Project-Id-Version: GNU hello-java 0.19-rc1\n"
"Report-Msgid-Bugs-To: bug-gnu-gettext@gnu.org\n"

#. Example: The set of prime numbers is {2, 3, 5, 7, 11, 13, ...}.
#: id=135956960462609535
msgid "The set of {$SET_NAME} is {{$XXX}, ...}."
msgstr ""

#: id=176798647517908084 pluralVar=EGGS_1
msgctxt "The number of eggs you need."
msgid "You have one egg"
msgid_plural "You have {$EGGS_2} eggs"
msgstr[0] "zYou zhave zone zegg"
msgstr[1] "zYou zhave zfew zeggs"
msgstr[2] "zYou zhave z{$EGGS_2} zeggs"

#: id=123
msgid ""
"ID Line 1\n"
"ID Line 2\n"
"ID Line 3"
msgstr ""
"STR Line 1\n"
"STR Line 2\n"
"STR Line 3"

`[1:]

var file = File{
	Header: textproto.MIMEHeader{
		"Content-Transfer-Encoding": {"8bit"},
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Language":                  {"sk"},
		"Language-Team":             {"Slovak <sk-i18n@lists.linux.sk>"},
		"Last-Translator":           {"Marcel Telka <marcel@telka.sk>"},
		"Mime-Version":              {"1.0"},
		"Plural-Forms":              {"nplurals=3; plural=(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2;"},
		"Po-Revision-Date":          {"2014-05-10 18:15+0200"},
		"Project-Id-Version":        {"GNU hello-java 0.19-rc1"},
		"Report-Msgid-Bugs-To":      {"bug-gnu-gettext@gnu.org"},
	},
	Messages: []Message{
		{
			Comment: Comment{
				ExtractedComments: []string{"Example: The set of prime numbers is {2, 3, 5, 7, 11, 13, ...}."},
				References:        []string{"id=135956960462609535"},
			},
			Id:  "The set of {$SET_NAME} is {{$XXX}, ...}.",
			Str: []string{""},
		},

		{
			Comment: Comment{
				ExtractedComments: nil,
				References:        []string{"id=176798647517908084", "pluralVar=EGGS_1"},
			},
			Ctxt:     "The number of eggs you need.",
			Id:       "You have one egg",
			IdPlural: "You have {$EGGS_2} eggs",
			Str: []string{
				"zYou zhave zone zegg",
				"zYou zhave zfew zeggs",
				"zYou zhave z{$EGGS_2} zeggs",
			},
		},

		{
			Comment: Comment{
				References: []string{"id=123"},
			},
			Id:  "ID Line 1\nID Line 2\nID Line 3",
			Str: []string{"STR Line 1\nSTR Line 2\nSTR Line 3"},
		},
	}}

func TestParse(t *testing.T) {
	var actual, err = Parse(strings.NewReader(po))
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(file.Header, actual.Header) {
		t.Errorf("expected header(%v):\n%v\ngot header(%v):\n%v",
			len(file.Header), file.Header, len(actual.Header), actual.Header)
	}
	if !reflect.DeepEqual(file.Messages, actual.Messages) {
		t.Errorf("expected msgs:\n%v\ngot msgs:\n%v", file.Messages, actual.Messages)
	}
}

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	var n, err = file.WriteTo(&buf)
	if err != nil {
		t.Error(err)
	}
	if n != int64(buf.Len()) {
		t.Errorf("n (%v) != buf length (%v)", n, buf.Len())
	}

	if buf.String() != po {
		t.Errorf("expected:\n%v\ngot:\n%v", po, buf.String())
	}

	actualLines := strings.Split(buf.String(), "\n")
	expLines := strings.Split(po, "\n")
	for i := range expLines {
		if expLines[i] != actualLines[i] {
			t.Errorf("%q != %q\n", expLines[i], actualLines[i])
		}
	}
}
