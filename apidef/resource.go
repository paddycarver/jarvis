package apidef

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// A Resource is a representation of a resource that can be manipulated through an API.
type Resource struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Parent       string        `json:"parent,omitempty"`
	URLSlug      string        `json:"url_slug"`
	Properties   []Property    `json:"properties"`
	Interactions []Interaction `json:"interactions,omitempty"`
}

// A Property is a definition of a specific field or property in a resource being returned by an API. It contains the information and constraints about the field.
type Property struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Format      string `json:"format,omitempty"`
	Maximum     int    `json:"maximum,omitempty"`
	Minimum     int    `json:"minimum,omitempty"`
	ValueType   string `json:"value_type,omitempty"`
}

// An Interaction is the definition of a specific action that can be performed against a resource using the API. It contains the information and constraints of that action.
type Interaction struct {
	ID                   string   `json:"id"`
	Verb                 string   `json:"verb"`
	Description          string   `json:"description"`
	OmittedInputFields   []string `json:"omitted_input_fields,omitempty"`
	RejectedInputFields  []string `json:"rejected_input_fields,omitempty"`
	RequiredInputFields  []string `json:"required_input_fields,omitempty"`
	OmittedOutputFields  []string `json:"omitted_output_fields,omitempty"`
	RejectedOutputFields []string `json:"rejected_output_fields,omitempty"`
	RequiredOutputFields []string `json:"required_output_fields,omitempty"`
	InputSource          string   `json:"input_source,omitempty"`
}

// ParseFile will read the specified resource file and parse it into a Resource, which is then returned.
func ParseFile(path string) (Resource, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Resource{}, err
	}
	var resource Resource
	err = json.Unmarshal(content, &resource)
	if err != nil {
		return Resource{}, err
	}
	return resource, nil
}

// Parse will find all resource files in the specified directory and parse them into Resources, which are then returned.
func Parse(path string) (map[string]Resource, error) {
	results := map[string]Resource{}
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err // if there's an error walking into the file or dir, fail
		}
		if info.IsDir() {
			return nil // skip sub-directories
		}
		r, err := ParseFile(p)
		if err != nil {
			return err
		}
		results[r.ID] = r
		return nil
	})
	return results, err
}
