package poeditor

import (
	"net/url"
	"strconv"
)

// Project represents a POEditor project
type Project struct {
	POEditor *POEditor
	ID       int
}

// ListLanguages lists all the available languages in the project
func (p *Project) ListLanguages() ([]Language, error) {
	res := languagesResult{}
	err := p.post("/languages/list", url.Values{}, &res)
	if err != nil {
		return []Language{}, err
	}
	ls := make([]Language, len(res.Languages))
	for i, l := range res.Languages {
		ls[i] = Language{Project: p, Code: l.Code}
	}
	return ls, nil
}

func (p *Project) post(endpoint string, params url.Values, res interface{}) error {
	params["id"] = []string{strconv.Itoa(p.ID)}
	return p.POEditor.post(endpoint, params, res)
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
