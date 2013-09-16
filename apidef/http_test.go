package apidef

import (
	"testing"
)

type endpointPieces struct {
	resource    *Resource
	interaction *Interaction
}

var (
	rootResource = &Resource{
		ID:        "rootResource",
		URLPrefix: "roots",
		URLSlug:   "id",
	}
	childResource = &Resource{
		ID:        "childResource",
		URLPrefix: "children",
		URLSlug:   "name",
		Parent:    rootResource,
	}
	orphanResource = &Resource{
		ID:                 "orphanResource",
		URLPrefix:          "orphans",
		URLSlug:            "birthday",
		Parent:             rootResource,
		ParentIsCollection: true,
	}
	grandchildResource = &Resource{
		ID:        "grandchildResource",
		URLPrefix: "grandchildren",
		URLSlug:   "id",
		Parent:    childResource,
	}
	orphanChildResource = &Resource{
		ID:        "orphanChildResource",
		URLPrefix: "orphanchildren",
		URLSlug:   "id",
		Parent:    orphanResource,
	}
	childOrphanResource = &Resource{
		ID:                 "childOrphanResource",
		URLPrefix:          "orphans",
		URLSlug:            "name",
		Parent:             childResource,
		ParentIsCollection: true,
	}
	orphanOrphanResource = &Resource{
		ID:                 "orphanOrphanResource",
		URLPrefix:          "orphans",
		URLSlug:            "name",
		Parent:             orphanResource,
		ParentIsCollection: true,
	}
)

var (
	listInteraction = &Interaction{
		ID:          "list",
		Name:        "list",
		Verb:        "list",
		Description: "list resources",
		AcceptMany:  false,
	}
	getInteraction = &Interaction{
		ID:          "get",
		Name:        "get",
		Verb:        "get",
		Description: "get resource",
		AcceptMany:  false,
	}
	updateInteraction = &Interaction{
		ID:          "update",
		Name:        "update",
		Verb:        "update",
		Description: "update resource",
		AcceptMany:  false,
	}
	createInteraction = &Interaction{
		ID:          "create",
		Name:        "create",
		Verb:        "create",
		Description: "create resource",
		AcceptMany:  false,
	}
	createManyInteraction = &Interaction{
		ID:          "createMany",
		Name:        "create many",
		Verb:        "create",
		Description: "create resources",
		AcceptMany:  true,
	}
	destroyInteraction = &Interaction{
		ID:          "destroy",
		Name:        "destroy",
		Verb:        "destroy",
		Description: "destroy resource",
		AcceptMany:  false,
	}
	destroyManyInteraction = &Interaction{
		ID:          "destroyMany",
		Name:        "destroy many",
		Verb:        "destroy",
		Description: "destroy resources",
		AcceptMany:  true,
	}
)

var testPaths = map[endpointPieces][]string{
	endpointPieces{rootResource, listInteraction}:         []string{"roots"},
	endpointPieces{childResource, listInteraction}:        []string{"roots", "{id}", "children"},
	endpointPieces{orphanResource, listInteraction}:       []string{"roots", "orphans"},
	endpointPieces{grandchildResource, listInteraction}:   []string{"roots", "{id}", "children", "{name}", "grandchildren"},
	endpointPieces{orphanChildResource, listInteraction}:  []string{"roots", "orphans", "{birthday}", "orphanchildren"},
	endpointPieces{childOrphanResource, listInteraction}:  []string{"roots", "{id}", "children", "orphans"},
	endpointPieces{orphanOrphanResource, listInteraction}: []string{"roots", "orphans", "orphans"},

	endpointPieces{rootResource, getInteraction}:         []string{"roots", "{id}"},
	endpointPieces{childResource, getInteraction}:        []string{"roots", "{id}", "children", "{name}"},
	endpointPieces{orphanResource, getInteraction}:       []string{"roots", "orphans", "{birthday}"},
	endpointPieces{grandchildResource, getInteraction}:   []string{"roots", "{id}", "children", "{name}", "grandchildren", "{id}"},
	endpointPieces{orphanChildResource, getInteraction}:  []string{"roots", "orphans", "{birthday}", "orphanchildren", "{id}"},
	endpointPieces{childOrphanResource, getInteraction}:  []string{"roots", "{id}", "children", "orphans", "{name}"},
	endpointPieces{orphanOrphanResource, getInteraction}: []string{"roots", "orphans", "orphans", "{name}"},

	endpointPieces{rootResource, updateInteraction}:         []string{"roots", "{id}"},
	endpointPieces{childResource, updateInteraction}:        []string{"roots", "{id}", "children", "{name}"},
	endpointPieces{orphanResource, updateInteraction}:       []string{"roots", "orphans", "{birthday}"},
	endpointPieces{grandchildResource, updateInteraction}:   []string{"roots", "{id}", "children", "{name}", "grandchildren", "{id}"},
	endpointPieces{orphanChildResource, updateInteraction}:  []string{"roots", "orphans", "{birthday}", "orphanchildren", "{id}"},
	endpointPieces{childOrphanResource, updateInteraction}:  []string{"roots", "{id}", "children", "orphans", "{name}"},
	endpointPieces{orphanOrphanResource, updateInteraction}: []string{"roots", "orphans", "orphans", "{name}"},

	endpointPieces{rootResource, createInteraction}:         []string{"roots"},
	endpointPieces{childResource, createInteraction}:        []string{"roots", "{id}", "children"},
	endpointPieces{orphanResource, createInteraction}:       []string{"roots", "orphans"},
	endpointPieces{grandchildResource, createInteraction}:   []string{"roots", "{id}", "children", "{name}", "grandchildren"},
	endpointPieces{orphanChildResource, createInteraction}:  []string{"roots", "orphans", "{birthday}", "orphanchildren"},
	endpointPieces{childOrphanResource, createInteraction}:  []string{"roots", "{id}", "children", "orphans"},
	endpointPieces{orphanOrphanResource, createInteraction}: []string{"roots", "orphans", "orphans"},

	endpointPieces{rootResource, createManyInteraction}:         []string{"roots"},
	endpointPieces{childResource, createManyInteraction}:        []string{"roots", "{id}", "children"},
	endpointPieces{orphanResource, createManyInteraction}:       []string{"roots", "orphans"},
	endpointPieces{grandchildResource, createManyInteraction}:   []string{"roots", "{id}", "children", "{name}", "grandchildren"},
	endpointPieces{orphanChildResource, createManyInteraction}:  []string{"roots", "orphans", "{birthday}", "orphanchildren"},
	endpointPieces{childOrphanResource, createManyInteraction}:  []string{"roots", "{id}", "children", "orphans"},
	endpointPieces{orphanOrphanResource, createManyInteraction}: []string{"roots", "orphans", "orphans"},

	endpointPieces{rootResource, destroyInteraction}:         []string{"roots", "{id}"},
	endpointPieces{childResource, destroyInteraction}:        []string{"roots", "{id}", "children", "{name}"},
	endpointPieces{orphanResource, destroyInteraction}:       []string{"roots", "orphans", "{birthday}"},
	endpointPieces{grandchildResource, destroyInteraction}:   []string{"roots", "{id}", "children", "{name}", "grandchildren", "{id}"},
	endpointPieces{orphanChildResource, destroyInteraction}:  []string{"roots", "orphans", "{birthday}", "orphanchildren", "{id}"},
	endpointPieces{childOrphanResource, destroyInteraction}:  []string{"roots", "{id}", "children", "orphans", "{name}"},
	endpointPieces{orphanOrphanResource, destroyInteraction}: []string{"roots", "orphans", "orphans", "{name}"},

	endpointPieces{rootResource, destroyManyInteraction}:         []string{"roots"},
	endpointPieces{childResource, destroyManyInteraction}:        []string{"roots", "{id}", "children"},
	endpointPieces{orphanResource, destroyManyInteraction}:       []string{"roots", "orphans"},
	endpointPieces{grandchildResource, destroyManyInteraction}:   []string{"roots", "{id}", "children", "{name}", "grandchildren"},
	endpointPieces{orphanChildResource, destroyManyInteraction}:  []string{"roots", "orphans", "{birthday}", "orphanchildren"},
	endpointPieces{childOrphanResource, destroyManyInteraction}:  []string{"roots", "{id}", "children", "orphans"},
	endpointPieces{orphanOrphanResource, destroyManyInteraction}: []string{"roots", "orphans", "orphans"},
}

func TestPathBuilding(t *testing.T) {
	for endpoint, pieces := range testPaths {
		pathPieces := endpoint.resource.BuildPathPieces(endpoint.interaction)
		if len(pathPieces) != len(pieces) {
			t.Errorf("Error building path for %s. Expected %d pieces in the path, got %d pieces.", endpoint.resource.ID+"#"+endpoint.interaction.ID, len(pathPieces), len(pieces))
		}
		for k, v := range pathPieces {
			if v != pieces[k] {
				t.Errorf("Error building path for %s, mismatch on piece %d. Expected %s, got %s.", endpoint.resource.ID+"#"+endpoint.interaction.ID, k, v, pieces[k])
			}
		}
	}
}
