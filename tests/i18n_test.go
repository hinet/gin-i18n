package tests

import (
	"golang.org/x/text/language"
	"hinet/gin-i18n/i18n"
	"log"
	"testing"
)

func test(t *testing.T) {
	translator := i18n.Translator{}
	err := translator.LoadLanguage()
	if err != nil {
		log.Fatalf("load i18n file failed %v", err)
	}
	translator.SetLanguage(language.Make("zh-CN"))
	translator.Translate("home.welcome!", nil)
}
