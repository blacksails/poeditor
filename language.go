package poeditor

import (
	"encoding/json"
	"io"
)

// Language represents a single language of a project
type Language struct {
	Project *Project
	Code    string
}

// AvailableLanguages lists all languages supported by POEditor. This is handy
// when you want to look up a particular language code.
func (poe *POEditor) AvailableLanguages() ([]AvailableLanguage, error) {
	var res []AvailableLanguage
	err := poe.post("/languages/available", nil, nil, &res)
	return res, err
}

// ListLanguages lists all the available languages in the project
func (p *Project) ListLanguages() ([]Language, error) {
	res := languagesResult{}
	err := p.post("/languages/list", nil, nil, &res)
	if err != nil {
		return []Language{}, err
	}
	ls := make([]Language, len(res.Languages))
	for i, l := range res.Languages {
		ls[i] = Language{Project: p, Code: l.Code}
	}
	return ls, nil
}

// AddLanguage adds a new language to the project. See
// POEditor.AvailableLanguages for a list of supported language codes.
func (p *Project) AddLanguage(code string) error {
	return p.post("/languages/add", map[string]string{"language": code}, nil, nil)
}

// Update inserts or overwrites translations for a language
// TODO: add fuzzy_trigger
func (l *Language) Update(terms []TranslationTerm) (CountResult, error) {
	var res CountResult
	// Typecheck translations
	for _, t := range terms {
		c := t.Translation.Content
		if _, ok := c.(string); ok {
			continue
		}
		if _, ok := c.(Plural); ok {
			continue
		}
		return res, ErrTranslationInvalid
	}
	// Encode and send translations
	ts, err := json.Marshal(terms)
	if err != nil {
		return res, err
	}
	err = l.post("/languages/update", map[string]string{"data": string(ts)}, nil, &res)
	return res, err
}

// Delete deletes the language
func (l *Language) Delete() error {
	return l.post("/languages/delete", nil, nil, nil)
}

// AvailableLanguage is a language supported by POEditor
type AvailableLanguage struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (l Language) post(endpoint string, fields map[string]string, files map[string]io.Reader, res interface{}) error {
	if fields == nil {
		fields = make(map[string]string)
	}
	fields["language"] = l.Code
	return l.Project.post(endpoint, fields, nil, res)
}

type languagesResult struct {
	Languages []language
}

type language struct {
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	Translations int     `json:"translations"`
	Percentage   float32 `json:"percentage"`
}
