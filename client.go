package poeditor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// POEditor is the main type used to interact with POEditor
type POEditor struct {
	apiToken string
}

// New returns a new POEditor given a POEditor API Token
func New(apiToken string) *POEditor {
	return &POEditor{apiToken: apiToken}
}

type Error struct {
	Status  string
	Code    string
	Message string
}

func (p Error) Error() string {
	return fmt.Sprintf("%s %s: %s", p.Status, p.Code, p.Message)
}

func (poe *POEditor) ListProjects() ([]Project, error) {
	res := projectsResult{}
	err := poe.post("/projects/list", url.Values{}, &res)
	if err != nil {
		return []Project{}, nil
	}
	ps := make([]Project, len(res.Projects))
	for i, p := range res.Projects {
		ps[i] = Project{POEditor: poe, ID: p.ID}
	}
	return ps, nil
}

type projectsResult struct {
	Projects []project `json:"projects"`
}

type project struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Public            int       `json:"public"`
	Open              int       `json:"open"`
	ReferenceLanguage string    `json:"reference_language"`
	Terms             int       `json:"terms"`
	Created           time.Time `json:"created"`
}

func (poe *POEditor) post(endpoint string, params url.Values, res interface{}) error {
	params["api_token"] = []string{poe.apiToken}
	resp, err := http.PostForm(fmt.Sprintf("https://api.poeditor.com/v2%s", endpoint), params)
	if err != nil {
		return err
	}
	poeRes := poEditorResponse{Result: res}
	json.NewDecoder(resp.Body).Decode(&poeRes)
	code, err := strconv.Atoi(poeRes.Response.Code)
	if err != nil {
		return err
	}
	if code-http.StatusOK > 100 {
		return poeRes.Response.ToError()
	}
	return nil
}

// Project represents a POEditor project
type Project struct {
	POEditor *POEditor
	ID       int
}

func (p *Project) post(endpoint string, params url.Values, res interface{}) error {
	params["id"] = []string{strconv.Itoa(p.ID)}
	return p.POEditor.post(endpoint, params, res)
}

// Project returns a Project with the given id
func (poe *POEditor) Project(id int) *Project {
	return &Project{POEditor: poe, ID: id}
}

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

type Language struct {
	Project *Project
	Code    string
}

func (l Language) post(endpoint string, params url.Values, res interface{}) error {
	params["language"] = []string{l.Code}
	return l.Project.post(endpoint, params, res)
}

const (
	FileFormatPO             = "po"
	FileFormatPOT            = "pot"
	FileFormatMO             = "mo"
	FileFormatXLS            = "xls"
	FileFormatCSV            = "csv"
	FileFormatRESW           = "resw"
	FileFormatRESX           = "resx"
	FileFormatAndriodStrings = "andriod_strings"
	FileFormatAppleStrings   = "apple_strings"
	FileFormatXLIFF          = "xliff"
	FileFormatProperties     = "properties"
	FileFormatKeyValueJSON   = "key_value_json"
	FileFormatJSON           = "json"
	FileFormatXMB            = "xmb"
	FileFormatXTB            = "xtb"
)

const (
	FilterTranslated   = "translated"
	FilterUntranslated = "untranslated"
	FilterFuzzy        = "fuzzy"
	FilterNotFuzzy     = "not_fuzzy"
	FilterAutomatic    = "automatic"
	FilterNotAutomatic = "not_automatic"
	FilterProofread    = "proofread"
	FilterNotProofread = "not_proofread"
)

func (l Language) Export(fileFormat string, filters []string, tags []string, dest io.Writer) error {
	params := url.Values{"type": {fileFormat}}
	if len(filters) > 0 {
		jsonFilters, err := json.Marshal(filters)
		if err != nil {
			return err
		}
		params["filters"] = []string{string(jsonFilters)}
	}
	if len(tags) > 0 {
		jsonTags, err := json.Marshal(tags)
		if err != nil {
			return err
		}
		params["tags"] = []string{string(jsonTags)}
	}
	exportRes := exportResult{}
	err := l.post("/projects/export", params, &exportRes)
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

type exportResult struct {
	URL string `json:"url"`
}

type poEditorResponse struct {
	Response response    `json:"response"`
	Result   interface{} `json:"result"`
}

type response struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r response) ToError() Error {
	return Error{Status: r.Status, Code: r.Code, Message: r.Message}
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
