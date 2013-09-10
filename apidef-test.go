package main

import (
	"fmt"
	"github.com/iron-io/jarvis/apidef"
	"os"
	"text/tabwriter"
)

func main() {
	resources, err := apidef.Parse("sample-resources/", "mq")
	if err != nil {
		panic(err)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for id, resource := range resources {
		endpoints := resource.BuildEndpoints()
		if len(endpoints) < 1 {
			continue
		}
		fmt.Fprintf(w, "\n%s (%s)\n", resource.Name, id)
		for _, endpoint := range endpoints {
			querystring := ""
			for _, param := range endpoint.Params {
				if param.Default != nil {
					continue
				}
				if querystring != "" {
					querystring += "&"
				}
				querystring += param.ID + "={" + param.Type + "}"
				if param.Repeated {
					querystring += "&" + param.ID + "={" + param.Type + "}"
					querystring += "&" + param.ID + "={" + param.Type + "}"
				}
			}
			if querystring != "" {
				querystring = "?" + querystring
			}
			fmt.Fprintf(w, "\t%s /%s\t%s\n", endpoint.Verb, endpoint.Path+querystring, endpoint.Name)
		}
	}
	w.Flush()
}
