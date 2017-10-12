package poeditor

import "fmt"

// Error represents an error sent from the POEditor API
type Error struct {
	Status  string
	Code    string
	Message string
}

func (p Error) Error() string {
	return fmt.Sprintf("%s %s: %s", p.Status, p.Code, p.Message)
}
