package i18n

import (
	"golang.org/x/text/language"
)

// Matcher 实现 Accept-Language 的解析
var Matcher = language.NewMatcher([]language.Tag{
	language.English,
	language.Chinese,
	language.Thai,
	// 添加其他语言...
})

// GetPreferredLanguage 解析 Accept-Language 头并返回最佳匹配的语言
func GetPreferredLanguage(acceptLanguage string) language.Tag {
	tag, _, _ := Matcher.Match(language.Make(acceptLanguage))
	return tag
}
