package poeditor

import "encoding/json"

// Term represents a POEditor Term. These are sent to the POEditor APIs using
// the Sync method.
type Term struct {
	Term        string       `json:"term"`
	Context     string       `json:"context,omitempty"`
	Plural      string       `json:"plural,omitempty"`
	Created     poEditorTime `json:"created,omitempty"`
	Updated     poEditorTime `json:"updated,omitempty"`
	Translation Translation  `json:"translation,omitempty"`
	Reference   string       `json:"reference,omitempty"`
	Comment     string       `json:"comment,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
}

// ListTerms returns all terms in the project
func (p *Project) ListTerms() ([]Term, error) {
	var res listTermsResult
	err := p.post("/terms/list", nil, nil, &res)
	return res.Terms, err
}

// ListTerms returns all terms in the project along with the translations for
// the language
func (l *Language) ListTerms() ([]Term, error) {
	var res listTermsResult
	err := l.post("/terms/list", nil, nil, &res)
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
