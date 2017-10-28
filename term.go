package poeditor

import (
	"encoding/json"
)

// Term represents a POEditor Term. These are sent to the POEditor APIs using
// the Sync method.
type Term struct {
	Term        string       `json:"term"`
	Context     string       `json:"context,omitempty"`
	Plural      string       `json:"plural,omitempty"`
	Created     poEditorTime `json:"created,omitempty"`
	Updated     poEditorTime `json:"updated,omitempty"`
	Translation interface{}  `json:"translation"`
	Reference   string       `json:"reference,omitempty"`
	Comment     string       `json:"comment,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
}

type _Term Term

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Term) UnmarshalJSON(bytes []byte) (err error) {
	var t2 _Term
	if err = json.Unmarshal(bytes, &t2); err != nil {
		return
	}
	*t = Term(t2)
	if s, ok := t.Translation.(string); ok {
		t.Translation = Singular(s)
	}
	if m, ok := t.Translation.(map[string]interface{}); ok {
		p := Plural{}
		if s, ok := m["one"].(string); ok {
			p.One = s
		}
		if s, ok := m["other"].(string); ok {
			p.Other = s
		}
		t.Translation = p
	}
	return nil
}
