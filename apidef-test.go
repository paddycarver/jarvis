package main

import (
	"fmt"
	"github.com/iron-io/jarvis/apidef"
)

func main() {
	resources, err := apidef.Parse("sample-resources/mq")
	if err != nil {
		panic(err)
	}
	for k, v := range resources {
		fmt.Printf("%s: %s\n", k, v.Description)
	}
}
