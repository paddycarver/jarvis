package spec

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/paddyforan/jarvis/parse"
	"math/big"
	"strings"
	"time"
)

// An Endpoint is a representation of an HTTP endpoint, the requests it expects, and the responses it can generate.
type Endpoint struct {
	Verb           string
	Path           string
	Params         []parse.Property
	Description    string
	Name           string
	SampleRequest  []byte
	SampleResponse []byte
}

func expectBody(verb string) bool {
	verb = strings.ToLower(verb)
	return verb == "create" || verb == "update"
}

func expectSlug(verb string) bool {
	verb = strings.ToLower(verb)
	return verb == "get" || verb == "update" || verb == "destroy"
}

func getHTTPVerb(verb string) string {
	verb = strings.ToLower(verb)
	switch verb {
	case "list":
		return "GET"
	case "get":
		return "GET"
	case "update":
		return "PUT"
	case "create":
		return "POST"
	case "destroy":
		return "DELETE"
	}
	return ""
}

// BuildEndpoints examines the resource it is called on and uses its properties to create and return a slice of endpoints.
func BuildEndpoints(r parse.Resource) ([]Endpoint, error) {
	endpoints := make([]Endpoint, len(r.Interactions))
	for i, interaction := range r.Interactions {
		if expectBody(interaction.Verb) {
			req, err := buildSampleRequest(r, &interaction)
			if err != nil {
				return endpoints, err
			}
			endpoints[i].SampleRequest = req
		}
		// TODO: generate sample response body and store it in endpoints[i].SampleResponse
		endpoints[i].Verb = getHTTPVerb(interaction.Verb)
		endpoints[i].Description = interaction.Description
		endpoints[i].Name = interaction.Name
		endpoints[i].Params = interaction.Params
		endpoints[i].Path = BuildPath(r, &interaction)
	}
	return endpoints, nil
}

// BuildPathPieces examines the resource it is called on (and that resource's parents) to create the pieces of the URL endpoint for the supplied interaction.
func BuildPathPieces(r parse.Resource, i *parse.Interaction) []string {
	var pieces []string
	if r.Parent != nil {
		pieces = append(pieces, BuildPathPieces(*r.Parent, nil)...)
		if !r.ParentIsCollection {
			pieces = append(pieces, "{"+r.Parent.URLSlug+"}")
		}
	}
	pieces = append(pieces, r.URLPrefix)
	if i == nil || !expectSlug(i.Verb) || i.AcceptMany {
		return pieces
	}
	pieces = append(pieces, "{"+r.URLSlug+"}")
	return pieces
}

// BuildPath examines the resource it is called on (and that resource's parents) to create the URL endpoint for the supplied interaction.
func BuildPath(r parse.Resource, i *parse.Interaction) string {
	return strings.Join(BuildPathPieces(r, i), "/")
}

func buildSampleRequest(r parse.Resource, i *parse.Interaction) ([]byte, error) {
	data := make([]byte, 0)
	if !expectBody(i.Verb) {
		return data, nil
	}
	if len(r.Properties) == 0 {
		return data, nil
	}
	resources := []map[string]interface{}{}
	num := 1
	if i.AcceptMany {
		num = 3
	}
	for iter := 0; iter < num; iter++ {
		resource := map[string]interface{}{} // ALL the maps!
		for _, property := range r.Properties {
			if !property.HasPerm("w") {
				continue // if we can't write the property, don't include it in the request
			}
			val, err := genRandomValue(&property)
			if err != nil {
				return data, err
			}
			if val != nil {
				resource[property.ID] = val
			}
		}
		if len(resource) == 0 {
			continue
		}
		resources = append(resources, resource)
	}
	if len(resources) == 0 {
		return data, nil
	}
	request := map[string]interface{}{} // This is so ugly.
	if i.AcceptMany {
		request[r.URLPrefix] = resources
	} else {
		request[r.ID] = resources[0]
	}
	return json.Marshal(request)
}

func genRandomValue(p *parse.Property) (interface{}, error) {
	if p.Default != nil {
		include, err := genRandomBool()
		if err != nil {
			return nil, err
		}
		if !include {
			return nil, nil
		}
		return p.Default, nil
	}
	if p.Values != nil {
		return pickRandomValue(p.Values)
	}
	p.Type = strings.ToLower(p.Type)
	switch p.Type {
	case "string":
		return genRandomString(p.Minimum, p.Maximum)
	case "bytes":
		return genRandomBytes(p.Minimum, p.Maximum)
	case "duration":
		return genRandomInt(p.Minimum, p.Maximum)
	case "datetime":
		return genRandomTime(p.Minimum, p.Maximum)
	case "int":
		return genRandomInt(p.Minimum, p.Maximum)
	case "Float":
		return genRandomFloat(p.Minimum, p.Maximum)
	case "boolean":
		return genRandomBool()
	}
	// TODO: throw error
	return nil, nil
}

func genRandomString(min, max int) (string, error) {
	b, err := genRandomBytes(min, max)
	if err != nil {
		return "", err
	}
	en := base64.StdEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)
	return string(d), nil
}

func genRandomBytes(min, max int) ([]byte, error) {
	chars, err := genRandomInt(min, max)
	if err != nil {
		return []byte{}, err
	}
	b := make([]byte, chars)
	rand.Read(b)
	return b, nil
}

func genRandomInt(min, max int) (int64, error) {
	if max == 0 && min == 0 {
		max = 32
	}
	if max-min == 0 {
		return int64(min), nil
	}
	bigMax := big.NewInt(int64(max))
	bigMin := big.NewInt(int64(min))
	diff := big.NewInt(0)
	diff = diff.Sub(bigMax, bigMin)
	i, err := rand.Int(rand.Reader, diff)
	i = i.Add(i, bigMin)
	return i.Int64(), err
}

func genRandomTime(min, max int) (time.Time, error) {
	// TODO
	var t time.Time
	return t, nil
}

func genRandomFloat(min, max int) (float64, error) {
	// TODO
	var f float64
	return f, nil
}

func genRandomBool() (bool, error) {
	i, err := genRandomInt(0, 2)
	return i == 1, err
}

func pickRandomValue(vals []interface{}) (interface{}, error) {
	i, err := genRandomInt(0, len(vals))
	if err != nil {
		return nil, err
	}
	return vals[i], nil
}
