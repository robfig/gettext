package po

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

var po = `
#. Example: The set of prime numbers is {2, 3, 5, 7, 11, 13, ...}.
#: id=135956960462609535
msgid "The set of {$SET_NAME} is {{$XXX}, ...}."
msgstr ""

#: id=176798647517908084 pluralVar=EGGS_1
msgctxt "The number of eggs you need."
msgid "You have one egg"
msgid_plural "You have {$EGGS_2} eggs"
msgstr ""

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

var file = File{[]Message{
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
		Str:      []string{""},
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
	if !reflect.DeepEqual(file, actual) {
		t.Errorf("expected:\n%v\ngot:\n%v", file, actual)
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
}
