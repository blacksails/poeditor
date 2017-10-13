package poeditor

import (
	"encoding/json"
	"io"
	"net/http"
)

// Language represents a single language of a project
type Language struct {
	Project *Project
	Code    string
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

func (l Language) post(endpoint string, fields map[string]string, files map[string]io.Reader, res interface{}) error {
	if fields == nil {
		fields = make(map[string]string)
	}
	fields["language"] = l.Code
	return l.Project.post(endpoint, fields, nil, res)
}

type exportResult struct {
	URL string `json:"url"`
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
