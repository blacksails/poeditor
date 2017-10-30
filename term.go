package poeditor

import "encoding/json"

type baseTerm struct {
	Term    string `json:"term"`
	Context string `json:"context,omitempty"`
}

// TranslationTerm is used when updating translations for a language
type TranslationTerm struct {
	baseTerm
	Translation Translation `json:"translation,omitempty"`
}

// Term is used when adding new terms, syncing or listing terms
type Term struct {
	baseTerm
	Plural    string       `json:"plural,omitempty"`
	Reference string       `json:"reference,omitempty"`
	Comment   string       `json:"comment,omitempty"`
	Tags      []string     `json:"tags,omitempty"`
	Created   poEditorTime `json:"created,omitempty"`
	Updated   poEditorTime `json:"updated,omitempty"`
}

// UpdateTerm is used when updating terms
type UpdateTerm struct {
	Term
	NewTerm    string `json:"new_term,omitempty"`
	NewContext string `json:"new_context,omitempty"`
}

// AddComment is used when adding a comment to a term
type AddComment struct {
	baseTerm
	Comment string `json:"comment"`
}

// TranslatedTerm is used when listing a project's terms along with translations
// for a language
type TranslatedTerm struct {
	Term
	Translation Translation `json:"translation"`
}

// ListTerms returns all terms in the project
func (p *Project) ListTerms() ([]Term, error) {
	var res listTermsResult
	err := p.post("/terms/list", nil, nil, &res)
	return res.Terms, err
}

// ListTerms returns all terms in the project along with the translations for
// the language
func (l *Language) ListTerms() ([]TranslatedTerm, error) {
	var res listTranslatedTermsResult
	err := l.post("/terms/list", nil, nil, &res)
	return res.Terms, err
}

// AddTerms adds the given terms to the project
func (p *Project) AddTerms(terms []Term) (CountResult, error) {
	var res termsCountResult
	jsonTerms, err := json.Marshal(terms)
	if err != nil {
		return res.Terms, err
	}
	err = p.post("/terms/add", map[string]string{"data": string(jsonTerms)}, nil, &res)
	return res.Terms, err
}

// UpdateTerms lets you change the text, context, reference, plural and tags of
// terms. Setting fuzzyTrigger to true marks associated translations as fuzzy.
func (p *Project) UpdateTerms(terms []UpdateTerm, fuzzyTrigger bool) (CountResult, error) {
	var res termsCountResult
	jsonTerms, err := json.Marshal(terms)
	if err != nil {
		return res.Terms, err
	}
	fuzzy := "0"
	if fuzzyTrigger {
		fuzzy = "1"
	}
	err = p.post("/terms/update", map[string]string{
		"data":          string(jsonTerms),
		"fuzzy_trigger": fuzzy,
	}, nil, &res)
	return res.Terms, err
}

// Translation is used to update translations in POEditor. The field Content
// must be either a string or a Plural type.
type Translation struct {
	Content   interface{}  `json:"content"`
	Fuzzy     int          `json:"fuzzy"`
	Proofread int          `json:"proofread"`
	Updated   poEditorTime `json:"updated"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Translation) UnmarshalJSON(bytes []byte) (err error) {
	type Alias Translation
	var t2 Alias
	if err = json.Unmarshal(bytes, &t2); err != nil {
		return
	}
	*t = Translation(t2)
	c := t.Content
	if s, ok := c.(string); ok {
		t.Content = string(s)
	}
	if m, ok := c.(map[string]interface{}); ok {
		p := Plural{}
		if s, ok := m["one"].(string); ok {
			p.One = s
		}
		if s, ok := m["other"].(string); ok {
			p.Other = s
		}
		t.Content = p
	}
	return nil
}

// Plural is a plural translation
type Plural struct {
	One   string `json:"one"`
	Other string `json:"other"`
}

type listTermsResult struct {
	Terms []Term
}

type listTranslatedTermsResult struct {
	Terms []TranslatedTerm
}
