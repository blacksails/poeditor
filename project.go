package poeditor

import (
	"encoding/json"
	"errors"
	"io"
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
	err := poe.post("/projects/view", map[string]string{"id": strconv.Itoa(id)}, nil, &res)
	if err != nil {
		return nil, err
	}
	return res.Project.toProject(poe), nil
}

// AddProject creates a new project with the given name and description
func (poe *POEditor) AddProject(name, description string) (*Project, error) {
	res := projectResult{}
	err := poe.post("/projects/add", map[string]string{
		"name":        name,
		"description": description,
	}, nil, &res)
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
	fields := make(map[string]string)
	for k, v := range props {
		if k == "name" || k == "description" || k == "reference_language" {
			fields[k] = v
		} else {
			return nil, ErrorUpdateFields
		}
	}
	err := p.post("/projects/update", fields, nil, &res)
	if err != nil {
		return nil, err
	}
	return res.Project.toProject(p.POEditor), nil
}

// Delete does its thing
func (p *Project) Delete() error {
	return p.post("/projects/delete", nil, nil, nil)
}

// Upload uploads terms
func (p *Project) Upload(reader io.Reader, options UploadOptions) (UploadResult, error) {
	var (
		res    UploadResult
		fields map[string]string
	)
	var validUpdating = func() bool {
		u := options.Updating
		return u == UploadTerms || u == UploadTranslations || u == UploadTermsTranslations
	}
	if !validUpdating() {
		return UploadResult{}, ErrorUploadUpdating
	}
	if options.Updating != UploadTerms && options.Language.Code == "" {
		return UploadResult{}, ErrorUploadLanguage
	}
	fields["language"] = options.Language.Code
	if options.Overwrite {
		fields["overwrite"] = "1"
	}
	if options.SyncTerms {
		fields["sync_terms"] = "1"
	}
	if len(options.Tags) > 0 {
		jsonTags, err := json.Marshal(options.Tags)
		if err != nil {
			return UploadResult{}, err
		}
		fields["tags"] = string(jsonTags)
	}
	if options.ReadFromSource {
		fields["read_from_source"] = "1"
	}
	if options.FuzzyTrigger {
		fields["fuzzy_trigger"] = "1"
	}
	err := p.post("/projects/upload", fields, map[string]io.Reader{"file": reader}, &res)
	if err != nil {
		return UploadResult{}, err
	}
	return res, nil
}

var (
	// ErrorUploadUpdating is returned from Upload when the value of Updating is invalid
	ErrorUploadUpdating = errors.New("Updating must be one of terms, terms_translations or translations")
	// ErrorUploadLanguage is return when language code is missing
	ErrorUploadLanguage = errors.New("Language code is required when uploading translations")
	// ErrorUpdateFields is returned when passing invalid fields to update
	ErrorUpdateFields = errors.New("Tried to update invalid field. Valid fields are name, description, reference_language")
)

const (
	// UploadTerms is a valid value of UploadOptions.Updating
	UploadTerms = "terms"
	// UploadTermsTranslations is a valid value of UploadOptions.Updating
	UploadTermsTranslations = "terms_translations"
	// UploadTranslations is a valid value of UploadOptions.Updating
	UploadTranslations = "translations"
)

// UploadOptions specifies options for upload of a file
type UploadOptions struct {
	Updating       string
	Language       Language
	Overwrite      bool
	SyncTerms      bool
	Tags           []string
	ReadFromSource bool
	FuzzyTrigger   bool
}

// UploadResult is returned when uploading a file
type UploadResult struct {
	Terms        UploadCountResult `json:"terms"`
	Translations UploadCountResult `json:"translations"`
}

// UploadCountResult is a part of UploadResult and shows counts for uploads
type UploadCountResult struct {
	Parsed  int `json:"parsed"`
	Added   int `json:"added"`
	Deleted int `json:"deleted"`
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

func (p *Project) post(endpoint string, fields map[string]string, files map[string]io.Reader, res interface{}) error {
	if fields == nil {
		fields = make(map[string]string)
	}
	fields["id"] = strconv.Itoa(p.ID)
	return p.POEditor.post(endpoint, fields, files, res)
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
