package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iron-io/jarvis/apidef"
)

func main() {
	resources, err := apidef.Parse("sample-resources/", "mq")
	if err != nil {
		panic(err)
	}
	for id, resource := range resources {
		endpoints, err := resource.BuildEndpoints()
		if err != nil {
			panic(err)
		}
		if len(endpoints) < 1 {
			continue
		}
		fmt.Printf("\n# %s (%s)\n", resource.Name, id)
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
			req := make([]byte, 0)
			if len(endpoint.SampleRequest) > 0 {
				var buf bytes.Buffer
				err := json.Indent(&buf, endpoint.SampleRequest, "\t", "  ")
				if err != nil {
					panic(err)
				}
        req = append(buf.Bytes(), []byte("\n")...)
			}
			fmt.Printf("## %s\n\n### Request\n\n%s /%s\n\n\t%s\n", endpoint.Name, endpoint.Verb, endpoint.Path+querystring, req)
		}
	}
}
