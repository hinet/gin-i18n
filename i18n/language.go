package i18n

import (
	"golang.org/x/text/language"
)

// Matcher Implement parsing of Accept-Language
var Matcher = language.NewMatcher([]language.Tag{
	language.English,
	language.Chinese,
	language.Thai,
	// Add other languages tag...
})

// GetPreferredLanguage Parse the Accept-Language header and return the best matching language
func GetPreferredLanguage(acceptLanguage string) language.Tag {
	tag, _, _ := Matcher.Match(language.Make(acceptLanguage))
	return tag
}
