package poeditor

// Term represents a POEditor Term. These are sent to the POEditor APIs using
// the Sync method.
type Term struct {
	Project   *Project `json:"-"`
	Term      string   `json:"term"`
	Context   string   `json:"context,omitempty"`
	Reference string   `json:"reference,omitempty"`
	Plural    string   `json:"plural,omitempty"`
	Comment   string   `json:"comment,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}
