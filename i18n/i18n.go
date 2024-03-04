package i18n

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

var (
	localLanguageFile string
	translator        *Translator
)

// MessageMap Store translation information for different languages
type MessageMap map[string]interface{}

func init() {
	flag.StringVar(&localLanguageFile, "lang", "config/lang/", "init language config")
}

// Translator translation provide
type Translator struct {
	translations map[language.Tag]MessageMap
	mu           sync.RWMutex
	language     language.Tag
}

// NewTranslator Create a new translator
func NewTranslator() *Translator {
	translator = &Translator{
		translations: make(map[language.Tag]MessageMap),
		language:     language.Make("en"),
	}
	if len(translator.translations) == 0 {
		err := translator.LoadLanguage()
		if err != nil {
			return translator
		}
	}
	return translator
}

// LoadLanguage Load the JSON language pack under the specified path
func (t *Translator) LoadLanguage() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	basePath := filepath.Join(filepath.Dir("."), localLanguageFile)
	files, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		filePath := filepath.Join(basePath, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
		var messages MessageMap
		if err := json.Unmarshal(content, &messages); err != nil {
			return err
		}
		filenameOutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		tag, _, _ := Matcher.Match(language.Make(filenameOutExt))
		t.translations[tag] = messages
	}

	return nil
}

func (t *Translator) SetLanguage(languageTag language.Tag) {
	t.language = languageTag
}

// Translate the given key based on language tags and replace variables
func (t *Translator) Translate(key string, params map[string]string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if messages, ok := t.translations[t.language]; ok {
		return t.recursive(key, messages, params)
	}
	return key // Default return key itself
}

// Recursive conversion of translated content
func (t *Translator) recursive(key string, messages MessageMap, params map[string]string) string {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return ""
	}
	currentKey := keys[0]
	translation, ok := messages[currentKey]
	if !ok {
		return currentKey
	}
	// If it is a string, replace the variable and return
	if str, ok := translation.(string); ok {
		if len(keys) == 1 {
			return interpolate(str, params)
		}
	}
	// If it is a nested map, continue recursive translation
	if nestedMap, ok := translation.(map[string]interface{}); ok {
		k := strings.Join(keys[1:], ".")
		return t.recursive(k, nestedMap, params)
	}
	return currentKey // If the value is not a string or a nested map, return the original key
}

// interpolate Replace variables in translation strings
func interpolate(translation string, params map[string]string) string {
	for key, value := range params {
		translation = strings.ReplaceAll(translation, "{"+key+"}", value)
	}
	return translation
}
//GetTranslator 
func GetTranslator() *Translator {
	if translator == nil {
		translator = NewTranslator()
		return translator
	}
	return translator
}
