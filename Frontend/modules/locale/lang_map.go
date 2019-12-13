package locale

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var languageMap = make(map[string]*po, 20)

func loadLanguages() {
	files, err := ioutil.ReadDir("./data/locales")
	if err != nil {
		fmt.Println("loadLanguages", err)
		return
	}
	for _, file := range files {
		if file.Name() == "templates.pot" || file.Name() == "." || file.Name() == ".." {
			continue
		}

		p, err := parseFile("./data/locales/" + file.Name())
		if err != nil {
			fmt.Println(file.Name(), ":", err)
			continue
		}
		if p == nil {
			fmt.Println(file.Name(), ":", "p is nil")
		}

		langName := strings.TrimPrefix(strings.TrimSuffix(file.Name(), ".po"), "templates-")
		languageMap[langName] = p
	}
}

func init() {
	loadLanguages()
}

// Get retrieves a string from a language
func Get(langs []string, str string, vars ...interface{}) string {
	for _, lang := range langs {
		l := languageMap[lang]

		if l == nil {
			continue
		}

		if el := l.Translations[str]; el != "" {
			if len(vars) > 0 {
				return fmt.Sprintf(el, vars...)
			}
			return el
		}
	}

	if len(vars) > 0 {
		return fmt.Sprintf(str, vars...)
	}

	return str
}
