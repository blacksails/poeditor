package poeditor

import (
	"encoding/json"
	"io"
	"net/http"
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

// ListProjects lists all the projects that are accessable by the used APIKey
func (poe *POEditor) ListProjects() ([]*Project, error) {
	res := projectsResult{}
	err := poe.post("/projects/list", nil, nil, &res)
	if err != nil {
		return []*Project{}, err
	}
	ps := make([]*Project, len(res.Projects))
	for i, p := range res.Projects {
		ps[i] = p.toProject(poe)
	}
	return ps, nil
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

// Sync syncs project terms with the given list of terms.
//
// CAUTION: this is a destructive operation. Any term not found in the input
// array will be deleted from the project.
func (p *Project) Sync(terms []Term) (CountResult, error) {
	jsonTerms, err := json.Marshal(terms)
	if err != nil {
		return CountResult{}, err
	}
	var res termsCountResult
	err = p.post("/projects/sync", map[string]string{"data": string(jsonTerms)}, nil, &res)
	if err != nil {
		return CountResult{}, err
	}
	return res.Terms, nil
}

// Export extracts the language in the given fileformat. For available file
// formats, see the FileFormat constants. Terms can be filtered using the
// Filter constants. Terms can also be filtered by tags.
func (l Language) Export(fileFormat string, filters []string, tags []string, dest io.Writer) error {
	fields := map[string]string{"type": fileFormat}
	if len(filters) > 0 {
		jsonFilters, err := json.Marshal(filters)
		if err != nil {
			return err
		}
		fields["filters"] = string(jsonFilters)
	}
	if len(tags) > 0 {
		jsonTags, err := json.Marshal(tags)
		if err != nil {
			return err
		}
		fields["tags"] = string(jsonTags)
	}
	exportRes := exportResult{}
	err := l.post("/projects/export", fields, nil, &exportRes)
	if err != nil {
		return err
	}
	export, err := http.Get(exportRes.URL)
	if err != nil {
		return err
	}
	_, err = io.Copy(dest, export.Body)
	return err
}

// ListTags returns a list of tags found on the project. This is not a standard
// API endpoint, but a useful helper never the less.
func (p *Project) ListTags() ([]string, error) {
	terms, err := p.ListTerms()
	if err != nil {
		return []string{}, err
	}
	tagsM := make(map[string]bool)
	for _, term := range terms {
		for _, tag := range term.Tags {
			tagsM[tag] = true
		}
	}
	tags := make([]string, len(tagsM))
	i := 0
	for tag := range tagsM {
		tags[i] = tag
	}
	return tags, nil
}

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

// CountResult is a part of UploadResult and returned directly from
// Project.Sync. It shows counts for uploaded/synced terms and translations
type CountResult struct {
	Parsed           int `json:"parsed"`
	Added            int `json:"added"`
	Deleted          int `json:"deleted"`
	WithAddedComment int `json:"with_added_comment"`
}

// UploadResult is returned when uploading a file
type UploadResult struct {
	Terms        CountResult `json:"terms"`
	Translations CountResult `json:"translations"`
}

const (
	// FileFormatPO specifies a .po file
	FileFormatPO = "po"
	// FileFormatPOT specifies a .pot file
	FileFormatPOT = "pot"
	// FileFormatMO specifies a .mo file
	FileFormatMO = "mo"
	// FileFormatXLS specifies an .xls file
	FileFormatXLS = "xls"
	// FileFormatCSV specifies a .csv file
	FileFormatCSV = "csv"
	// FileFormatRESW specifies an .resw file
	FileFormatRESW = "resw"
	// FileFormatRESX specifies an .resx file
	FileFormatRESX = "resx"
	// FileFormatAndroidStrings specifies strings should be in android format
	FileFormatAndroidStrings = "android_strings"
	// FileFormatAppleStrings specifies strings should be in apple format
	FileFormatAppleStrings = "apple_strings"
	// FileFormatXLIFF specifies an .xliff file
	FileFormatXLIFF = "xliff"
	// FileFormatProperties specifies a .propterties file
	FileFormatProperties = "properties"
	// FileFormatKeyValueJSON specifies a .json file in key value format
	FileFormatKeyValueJSON = "key_value_json"
	// FileFormatJSON specifies a .json file
	FileFormatJSON = "json"
	// FileFormatXMB specifies an .xmb file
	FileFormatXMB = "xmb"
	// FileFormatXTB specifies an .xtb file
	FileFormatXTB = "xtb"
)

const (
	// FilterTranslated filters terms in translated state
	FilterTranslated = "translated"
	// FilterUntranslated filters terms in untranslated state
	FilterUntranslated = "untranslated"
	// FilterFuzzy filters terms in fuzzy state
	FilterFuzzy = "fuzzy"
	// FilterNotFuzzy filters terms in not fuzzy state
	FilterNotFuzzy = "not_fuzzy"
	// FilterAutomatic filters terms in automatic state
	FilterAutomatic = "automatic"
	// FilterNotAutomatic filters terms in not automatic state
	FilterNotAutomatic = "not_automatic"
	// FilterProofread filters terms in proofread state
	FilterProofread = "proofread"
	// FilterNotProofread filters terms in not proofread state
	FilterNotProofread = "not_proofread"
)

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

type termsCountResult struct {
	Terms CountResult `json:"terms"`
}

type exportResult struct {
	URL string `json:"url"`
}
