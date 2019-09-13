package po

import "strings"

// PluralSelector returns the appropriate plural case to use, given a quantity.
type PluralSelector func(n int) int

var langNames = map[string]string{
	"ja":    "Japanese",
	"vi":    "Vietnamese",
	"ko":    "Korean",
	"zh":    "Chinese",
	"en":    "English",
	"de":    "German",
	"nl":    "Dutch",
	"sv":    "Swedish",
	"da":    "Danish",
	"no":    "Norwegian",
	"nb":    "Norwegian Bokmal",
	"nn":    "Norwegian Nynorsk",
	"fo":    "Faroese",
	"es":    "Spanish",
	"pt":    "Portuguese",
	"it":    "Italian",
	"bg":    "Bulgarian",
	"el":    "Greek",
	"fi":    "Finnish",
	"et":    "Estonian",
	"he":    "Hebrew",
	"eo":    "Esperanto",
	"hu":    "Hungarian",
	"tr":    "Turkish",
	"pt_BR": "Brazilian",
	"fr":    "French",
	"lv":    "Latvian",
	"ga":    "Irish",
	"ro":    "Romanian",
	"lt":    "Lithuanian",
	"ru":    "Russian",
	"uk":    "Ukrainian",
	"be":    "Belarusian",
	"sr":    "Serbian",
	"hr":    "Croatian",
	"cs":    "Czech",
	"sk":    "Slovak",
	"pl":    "Polish",
	"sl":    "Slovenian",
	"ar":    "Arabic",
	"ms":    "Malay",
}

// TODO: Fall back to these if Plural-Forms is not specified.
var pluralExprs = map[string]string{
	"ja":    "nplurals=1; plural=0;",
	"vi":    "nplurals=1; plural=0;",
	"ko":    "nplurals=1; plural=0;",
	"zh":    "nplurals=1; plural=0;",
	"ms":    "nplurals=1; plural=0;",
	"en":    "nplurals=2; plural=(n != 1);",
	"de":    "nplurals=2; plural=(n != 1);",
	"nl":    "nplurals=2; plural=(n != 1);",
	"sv":    "nplurals=2; plural=(n != 1);",
	"da":    "nplurals=2; plural=(n != 1);",
	"no":    "nplurals=2; plural=(n != 1);",
	"nb":    "nplurals=2; plural=(n != 1);",
	"nn":    "nplurals=2; plural=(n != 1);",
	"fo":    "nplurals=2; plural=(n != 1);",
	"es":    "nplurals=2; plural=(n != 1);",
	"pt":    "nplurals=2; plural=(n != 1);",
	"it":    "nplurals=2; plural=(n != 1);",
	"bg":    "nplurals=2; plural=(n != 1);",
	"el":    "nplurals=2; plural=(n != 1);",
	"fi":    "nplurals=2; plural=(n != 1);",
	"et":    "nplurals=2; plural=(n != 1);",
	"he":    "nplurals=2; plural=(n != 1);",
	"eo":    "nplurals=2; plural=(n != 1);",
	"hu":    "nplurals=2; plural=(n != 1);",
	"tr":    "nplurals=2; plural=(n != 1);",
	"pt_BR": "nplurals=2; plural=(n > 1);",
	"fr":    "nplurals=2; plural=(n > 1);",
	"lv":    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n != 0 ? 1 : 2);",
	"ga":    "nplurals=3; plural=n==1 ? 0 : n==2 ? 1 : 2;",
	"ro":    "nplurals=3; plural=n==1 ? 0 : (n==0 || (n%100 > 0 && n%100 < 20)) ? 1 : 2;",
	"lt":    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && (n%100<10 || n%100>=20) ? 1 : 2);",
	"ru":    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);",
	"uk":    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);",
	"be":    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);",
	"sr":    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);",
	"hr":    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);",
	"cs":    "nplurals=3; plural=(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2;",
	"sk":    "nplurals=3; plural=(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2;",
	"pl":    "nplurals=3; plural=(n==1 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);",
	"sl":    "nplurals=4; plural=(n%100==1 ? 0 : n%100==2 ? 1 : n%100==3 || n%100==4 ? 2 : 3);",
	"ar":    "nplurals=6; plural=(n==0 ? 0 : n==1 ? 1 : n==2 ? 2 : n%100>=3 && n%100<=10 ? 3 : n%100>=11 ? 4 : 5);",
}

// pluralSelectors contains a lookup from space-stripped plural forms strings to
// the functions that implement them.
var pluralSelectors = stripSpace(map[string]PluralSelector{
	"nplurals=1; plural=0;":                                                                                  plural0,
	"nplurals=2; plural=(n != 1);":                                                                           pluralNeq1,
	"nplurals=2; plural=(n > 1);":                                                                            pluralGt1,
	"nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n != 0 ? 1 : 2);":                                        pluralLatvian,
	"nplurals=3; plural=n==1 ? 0 : n==2 ? 1 : 2;":                                                            pluralIrish,
	"nplurals=3; plural=n==1 ? 0 : (n==0 || (n%100 > 0 && n%100 < 20)) ? 1 : 2;":                             pluralRomanian,
	"nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && (n%100<10 || n%100>=20) ? 1 : 2);":            pluralLithuanian,
	"nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);": pluralRussian,
	"nplurals=3; plural=(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2;":                                                pluralCzech,
	"nplurals=3; plural=(n==1 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);":                 pluralPolish,
	"nplurals=4; plural=(n%100==1 ? 0 : n%100==2 ? 1 : n%100==3 || n%100==4 ? 2 : 3);":                       pluralSlovenian,
	"nplurals=6; plural=(n==0 ? 0 : n==1 ? 1 : n==2 ? 2 : n%100>=3 && n%100<=10 ? 3 : n%100>=11 ? 4 : 5);":   pluralArabic,
})

func stripSpace(m map[string]PluralSelector) map[string]PluralSelector {
	var r = make(map[string]PluralSelector, len(m))
	for k, v := range m {
		r[strings.Replace(k, " ", "", -1)] = v
	}
	return r
}

// lookupPluralSelectors looks up the given plural form from the set of known ones.
// nil is returned if the plural form was not recognized.
func lookupPluralSelector(pluralForms string) PluralSelector {
	return pluralSelectors[strings.Replace(pluralForms, " ", "", -1)]
}

// PluralSelectorForLanguage returns the appropriate plural selector for the
// provided languge code. The code can be either the too letter code ("en") or
// the 5 character variant ("en_GB")
func PluralSelectorForLanguage(lang string) PluralSelector {
	lang = strings.Replace(lang, "-", "_", -1)
	if pluralForms, found := pluralExprs[lang]; found {
		return lookupPluralSelector(pluralForms)
	}
	if len(lang) > 2 && lang[2] == '_' {
		// Naively trim the input
		if pluralForms, found := pluralExprs[lang[:2]]; found {
			return lookupPluralSelector(pluralForms)
		}
	}
	return nil
}

func plural0(n int) int {
	return 0
}

func pluralNeq1(n int) int {
	if n != 1 {
		return 1
	}
	return 0
}

func pluralGt1(n int) int {
	if n > 1 {
		return 1
	}
	return 0
}

func pluralLatvian(n int) int {
	switch {
	case n%10 == 1 && n%100 != 11:
		return 0
	case n != 0:
		return 1
	default:
		return 2
	}
}

func pluralIrish(n int) int {
	switch n {
	case 1:
		return 0
	case 2:
		return 1
	default:
		return 2
	}
}

func pluralRomanian(n int) int {
	switch {
	case n == 1:
		return 0
	case n == 0 || (n%100 > 0 && n%100 < 20):
		return 1
	default:
		return 2
	}
}

func pluralLithuanian(n int) int {
	switch {
	case n%10 == 1 && n%100 != 11:
		return 0
	case n%10 >= 2 && (n%100 < 10 || n%100 >= 20):
		return 1
	default:
		return 2
	}
}

func pluralRussian(n int) int {
	switch {
	case n%10 == 1 && n%100 != 11:
		return 0
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20):
		return 1
	default:
		return 2
	}
}

func pluralCzech(n int) int {
	switch {
	case n == 1:
		return 0
	case n >= 2 && n <= 4:
		return 1
	default:
		return 2
	}
}

func pluralPolish(n int) int {
	switch {
	case n == 1:
		return 0
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20):
		return 1
	default:
		return 2
	}
}

func pluralSlovenian(n int) int {
	switch {
	case n%100 == 1:
		return 0
	case n%100 == 2:
		return 1
	case n%100 == 3 || n%100 == 4:
		return 2
	default:
		return 3
	}
}

func pluralArabic(n int) int {
	switch {
	case n == 0:
		return 0
	case n == 1:
		return 1
	case n == 2:
		return 2
	case n%100 >= 3 && n%100 <= 10:
		return 3
	case n%100 >= 11:
		return 4
	default:
		return 5
	}
}
