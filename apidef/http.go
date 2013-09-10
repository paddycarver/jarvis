package apidef

import (
	"strings"
)

// An Endpoint is a representation of an HTTP endpoint, the requests it expects, and the responses it can generate.
type Endpoint struct {
	Verb           string
	Path           string
	Params         []Property
	Description    string
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
func (r Resource) BuildEndpoints() []Endpoint {
	endpoints := make([]Endpoint, len(r.Interactions))
	for i, interaction := range r.Interactions {
		if expectBody(interaction.Verb) {
			// TODO: generate sample body and store it in endpoints[i].SampleRequest
		}
		// TODO: generate sample response body and store it in endpoints[i].SampleResponse
		endpoints[i].Verb = getHTTPVerb(interaction.Verb)
		endpoints[i].Description = interaction.Description
		endpoints[i].Params = interaction.Params
		endpoints[i].Path = r.BuildPath(&interaction)
	}
	return endpoints
}

// BuildPathPieces examines the resource it is called on (and that resource's parents) to create the pieces of the URL endpoint for the supplied interaction.
func (r Resource) BuildPathPieces(i *Interaction) []string {
	var pieces []string
	if r.Parent != nil {
		pieces = append(pieces, r.Parent.BuildPathPieces(nil)...)
		if !r.ParentIsCollection {
			pieces = append(pieces, "{"+r.Parent.URLSlug+"}")
		}
	}
	pieces = append(pieces, r.URLPrefix)
	if i == nil || !expectSlug(i.Verb) {
		return pieces
	}
	pieces = append(pieces, "{"+r.URLSlug+"}")
	return pieces
}

// BuildPath examines the resource it is called on (and that resource's parents) to create the URL endpoint for the supplied interaction.
func (r Resource) BuildPath(i *Interaction) string {
	return strings.Join(r.BuildPathPieces(i), "/")
}
