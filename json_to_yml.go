package main

import (
	"flag"
	"github.com/iron-io/jarvis/apidef"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"strings"
)

func main() {
	var file string
	flag.StringVar(&file, "f", "", "Path to file to convert.")
	flag.Parse()
	if file == "" {
		flag.Usage()
		return
	}
	resource, err := apidef.ParseFile(file)
	if err != nil {
		panic(err)
	}
	data, err := goyaml.Marshal(&resource)
	if err != nil {
		panic(err)
	}
	fileInfo, err := os.Lstat(file)
	if err != nil {
		panic(err)
	}
	file = strings.Replace(file, ".json", ".yml", -1)
	err = ioutil.WriteFile(file, data, fileInfo.Mode())
	if err != nil {
		panic(err)
	}
}
