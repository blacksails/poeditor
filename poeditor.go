package poeditor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// Project returns a Project with the given id
func (poe *POEditor) Project(id int) *Project {
	return &Project{POEditor: poe, ID: id}
}

func (poe *POEditor) post(endpoint string, fields map[string]string, files map[string]io.Reader, res interface{}) error {
	// Initiate fields if nil
	if fields == nil {
		fields = make(map[string]string)
	}
	// Set API Token
	fields["api_token"] = poe.apiToken
	// Initiate multipart writer
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	// Write key value fields
	for k, v := range fields {
		err := writer.WriteField(k, v)
		if err != nil {
			return err
		}
	}
	// Write files
	for k, v := range files {
		w, err := writer.CreateFormFile(k, k)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, v)
		if err != nil {
			return err
		}
	}
	err := writer.Close()
	if err != nil {
		return err
	}
	// Send request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.poeditor.com/v2%s", endpoint), &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// Decode response
	poeRes := poEditorResponse{Result: res}
	if os.Getenv("DEBUG") == "true" {
		var body bytes.Buffer
		json.NewDecoder(io.TeeReader(resp.Body, &body)).Decode(&poeRes)
		log.Println(body.String())
	} else {
		json.NewDecoder(resp.Body).Decode(&poeRes)
	}
	code, err := strconv.Atoi(poeRes.Response.Code)
	if err != nil {
		return err
	}
	if code-http.StatusOK > 100 {
		return poeRes.Response.ToError()
	}
	return nil
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

const poEditorTimeLayout string = "2006-01-02T15:04:05Z0700"

type poEditorTime struct {
	time.Time
}

func (t *poEditorTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	pt, err := time.Parse(poEditorTimeLayout, s)
	t.Time = pt
	return err
}
