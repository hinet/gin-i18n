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

// MessageMap 存储不同语言的翻译信息
type MessageMap map[string]interface{}

func init() {
	flag.StringVar(&localLanguageFile, "lang", "config/lang/", "init language config")
}

// Translator 提供翻译功能
type Translator struct {
	translations map[language.Tag]MessageMap
	mu           sync.RWMutex
	language     language.Tag
}

// NewTranslator 创建新的翻译器
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

// LoadLanguage 加载指定路径下的 JSON 语言包
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

// Translate 根据语言标签翻译给定的键，并进行变量替换
func (t *Translator) Translate(key string, params map[string]string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if messages, ok := t.translations[t.language]; ok {
		return t.recursively(key, messages, params)
	}
	return key // 默认返回键本身
}

// 递归转换翻译内容
func (t *Translator) recursively(key string, messages MessageMap, params map[string]string) string {
	// 如果已经到达最后一级键，返回对应的翻译结果
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return ""
	}
	// 获取当前级别的键
	currentKey := keys[0]
	// 获取当前级别的值
	translation, ok := messages[currentKey]
	if !ok {
		return currentKey // 如果没有找到对应的键，返回原始键
	}
	// 如果是字符串，进行变量替换并返回
	if str, ok := translation.(string); ok {
		// 如果已经是最后一级键，进行变量替换并返回
		if len(keys) == 1 {
			return interpolate(str, params)
		}
	}
	// 如果是嵌套的 map，则继续递归翻译
	if nestedMap, ok := translation.(map[string]interface{}); ok {
		k := strings.Join(keys[1:], ".")
		return t.recursively(k, nestedMap, params)
	}
	return currentKey // 如果值不是字符串也不是嵌套 map，则返回原始键
}

// interpolate 替换翻译字符串中的变量
func interpolate(translation string, params map[string]string) string {
	for key, value := range params {
		translation = strings.ReplaceAll(translation, "{"+key+"}", value)
	}
	return translation
}

func GetTranslator() *Translator {
	if translator == nil {
		translator = NewTranslator()
		return translator
	}
	return translator
}
