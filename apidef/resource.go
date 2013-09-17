package apidef

import (
	"errors"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"strings"
)

// A Resource is a representation of a resource that can be manipulated through an API.
type Resource struct {
	ID                 string        `yaml:"id"`
	Name               string        `yaml:"name"`
	Description        string        `yaml:"description"`
	Parent             *Resource     `yaml:"-"`
	ParentString       string        `yaml:"parent,omitempty"`
	ParentIsCollection bool          `yaml:"parent_is_collection,omitempty"`
	URLSlug            string        `yaml:"url_slug"`
	URLPrefix          string        `yaml:"url_prefix"`
	Properties         []Property    `yaml:"properties"`
	Interactions       []Interaction `yaml:"interactions,omitempty"`
}

// A Property is a definition of a specific field or property in a resource being returned by an API. It contains the information and constraints about the field.
type Property struct {
	ID          string        `yaml:"id"`
	Type        string        `yaml:"type"`
	Description string        `yaml:"description"`
	Values      []interface{} `yaml:"values,omitempty"`  // A list of acceptable values
	Default     interface{}   `yaml:"default,omitempty"` // The default value, if this property is optional
	Maximum     int           `yaml:"maximum,omitempty"`
	Minimum     int           `yaml:"minimum,omitempty"`
	Permissions []string      `yaml:"permissions,omitempty"` // Permissions clients have for this property. Acceptable values: r, w
	Repeated    bool          `yaml:"repeated,omitempty"`    // If this property can appear more than once in URL parameters
}

// An Interaction is the definition of a specific action that can be performed against a resource using the API. It contains the information and constraints of that action.
type Interaction struct {
	ID          string     `yaml:"id"`
	Name        string     `yaml:"name"`
	Verb        string     `yaml:"verb"`
	Description string     `yaml:"description"`
	Params      []Property `yaml:"params,omitempty"`      // Properties passed as URL params
	AcceptMany  bool       `yaml:"accept_many,omitempty"` // expect an array, not a single resource
}

// ParseFile will read the specified resource file and parse it into a Resource, which is then returned.
func ParseFile(path string) (Resource, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Resource{}, err
	}
	var resource Resource
	err = goyaml.Unmarshal(content, &resource)
	if err != nil {
		return Resource{}, err
	}
	return resource, nil
}

func createImportList(root, path string) (map[string]string, error) {
	results := map[string]string{}
	err := filepath.Walk(root+path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err // if there's an error walking into the file or dir, fail
		}
		if info.IsDir() {
			return nil // skip sub-directories
		}
		if !strings.HasSuffix(p, ".yml") {
			return nil // skip non-yaml files
		}
		results[strings.TrimPrefix(p, root)] = p
		return nil
	})
	return results, err
}

func importParents(root string, r Resource, cache map[string]bool) (map[string]*Resource, error) {
	if r.ParentString == "" {
		return map[string]*Resource{}, nil
	}
	path := getParentPath(r)
	if path == "" {
		return map[string]*Resource{}, nil
	}
	if _, ok := cache[path]; ok {
		return map[string]*Resource{}, nil
	}
	results, err := Parse(root, path)
	if err != nil {
		return map[string]*Resource{}, errors.New("Error parsing import " + path + ": " + err.Error())
	}
	cache[path] = true
	return results, nil
}

// Parse will find all resource files in the specified directory and parse them into Resources, which are then returned.
func Parse(root, path string) (map[string]*Resource, error) {
	results := map[string]*Resource{}
	importCache := map[string]bool{}

	toImport, err := createImportList(root, path)
	if err != nil {
		return results, err
	}

	for id, filePath := range toImport {
		r, err := ParseFile(filePath)
		if err != nil {
			return results, errors.New("Error parsing " + id + ": " + err.Error())
		}
		myPath := getResourcePath(id)
		if getParentPath(r) != myPath {
			parents, err := importParents(root, r, importCache)
			if err != nil {
				return results, err
			}
			for parentPath, parent := range parents {
				results[parentPath] = parent
			}
		}

		results[path+"/"+r.ID] = &r
	}

	// map our resources to their parents
	for k, r := range results {
		if r.ParentString == "" {
			continue
		}
		if p, ok := results[r.ParentString]; ok {
			results[k].Parent = p
		} else {
			return results, errors.New("Error parsing " + path + ": Parent of " + k + " not found: " + r.ParentString)
		}
	}
	return results, err
}

func getParentPath(resource Resource) string {
	if resource.ParentString == "" {
		return ""
	}
	return getResourcePath(resource.ParentString)
}

func getResourcePath(path string) string {
	index := strings.LastIndex(path, "/")
	if index == -1 {
		return ""
	}
	return path[0:index]
}

// HasPerm is a helper function that tests whether the passed in permission is present in the property's permissions list.
func (p Property) HasPerm(perm string) bool {
	for _, can := range p.Permissions {
		if strings.ToLower(perm) == strings.ToLower(can) {
			return true
		}
	}
	return false
}
