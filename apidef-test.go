package main

import (
	"fmt"
	"github.com/iron-io/jarvis/apidef"
)

func main() {
	resources, err := apidef.Parse("sample-resources/", "mq")
	if err != nil {
		panic(err)
	}
	for id, resource := range resources {
    endpoints := resource.BuildEndpoints()
    if len(endpoints) < 1 {
      continue
    }
		fmt.Printf("\n%s (%s)\n", resource.Name, id)
    for _, endpoint := range endpoints {
      fmt.Printf("\t%s /%s\t%s\n", endpoint.Verb, endpoint.Path, endpoint.Description)
    }
	}
}
