package poeditor

import (
	"errors"
	"fmt"
)

// Error represents an error sent from the POEditor API
type Error struct {
	Status  string
	Code    string
	Message string
}

func (p Error) Error() string {
	return fmt.Sprintf("%s %s: %s", p.Status, p.Code, p.Message)
}

var (
	// ErrorUploadUpdating is returned from Project.Upload when the value of
	// Updating is invalid
	ErrorUploadUpdating = errors.New("Updating must be one of terms, terms_translations or translations")
	// ErrorUploadLanguage is return when language code is missing
	ErrorUploadLanguage = errors.New("Language code is required when uploading translations")
 // ErrorUpdateFields is returned when passing invalid fields to Project.Update
	ErrorUpdateFields = errors.New("Tried to update invalid field. Valid fields are name, description, reference_language")
)
