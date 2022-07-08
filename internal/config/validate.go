package config

import (
	"embed"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/*
var schemas embed.FS

func validateConfigBytes(data []byte, version int) bool {
	if version != 1 {
		fmt.Fprintf(os.Stderr, "Config version %d is unsupported. Only version 1 is supported.\n", version)
		return false
	}
	body, err := schemas.ReadFile(fmt.Sprintf("schemas/schema-%d.yaml", version))
	if err != nil {
		panic(err)
	}
	schemaLoader, err := newRawLoaderFromYAMLBytes(body)
	if err != nil {
		panic(err)
	}
	docLoader, err := newRawLoaderFromYAMLBytes(data)
	if err != nil {
		panic(err)
	}
	result, err := gojsonschema.Validate(*schemaLoader, *docLoader)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if !result.Valid() {
		fmt.Fprintln(os.Stderr, "The config is not valid:")
		for _, desc := range result.Errors() {
			fmt.Fprintf(os.Stderr, " - %s\n", desc)
		}
		return false
	}
	return true
}

func newRawLoaderFromYAMLBytes(data []byte) (*gojsonschema.JSONLoader, error) {
	var document map[string]interface{}
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, err
	}
	loader := gojsonschema.NewRawLoader(document)
	return &loader, nil
}
