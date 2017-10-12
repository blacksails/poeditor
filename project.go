package poeditor

import (
	"net/url"
	"strconv"
	"time"
)

// Project represents a POEditor project
type Project struct {
	POEditor          *POEditor
	ID                int
	Name              string
	Description       string
	Public            int
	Open              int
	ReferenceLanguage string
	Terms             int
	Created           time.Time
}

// ViewProject returns project with given ID
func (poe *POEditor) ViewProject(id int) (*Project, error) {
	res := projectResult{}
	err := poe.post("/projects/view", url.Values{"id": []string{strconv.Itoa(id)}}, &res)
	if err != nil {
		return nil, err
	}
	return res.Project.toProject(poe), nil
}

// AddProject creates a new project with the given name and description
func (poe *POEditor) AddProject(name, description string) (*Project, error) {
	res := projectResult{}
	err := poe.post("/projects/add", url.Values{
		"name":        []string{name},
		"description": []string{description},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.Project.toProject(poe), nil
}

/*
Update updates the project according to the map[string]string.

		...
		p, err := p.Update(map[string]string{
			"name": "a project name",
			"description": "a project description"
			"reference_language": "a reference language code"
		})

Omitted key value pairs are not updated. Only `name`, `description` and
`reference_language` can be updated.
*/
func (p *Project) Update(props map[string]string) (*Project, error) {
	res := projectResult{}
	updates := url.Values{}
	for k, v := range props {
		if k == "name" || k == "description" || k == "reference_language" {
			updates[k] = []string{v}
		}
	}
	err := p.post("/projects/update", updates, &res)
	if err != nil {
		return nil, err
	}
	return res.Project.toProject(p.POEditor), nil
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

type projectsResult struct {
	Projects []project `json:"projects"`
}

type projectResult struct {
	Project project `json:"project"`
}

type project struct {
	ID                int          `json:"id"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Public            int          `json:"public"`
	Open              int          `json:"open"`
	ReferenceLanguage string       `json:"reference_language"`
	Terms             int          `json:"terms"`
	Created           poEditorTime `json:"created"`
}

func (p project) toProject(poe *POEditor) *Project {
	return &Project{
		POEditor:          poe,
		ID:                p.ID,
		Name:              p.Name,
		Description:       p.Description,
		Public:            p.Public,
		Open:              p.Open,
		ReferenceLanguage: p.ReferenceLanguage,
		Terms:             p.Terms,
		Created:           p.Created.Time,
	}
}
