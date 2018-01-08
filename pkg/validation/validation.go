package validation

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

// https://github.com/xeipuuv/gojsonschema#formats
type ValidFormat struct{}

// IsFormat always returns true and meets the
// gojsonschema.FormatChecker interface
func (f ValidFormat) IsFormat(input interface{}) bool {
	return true
}

// Based on https://stackoverflow.com/questions/40737122/convert-yaml-to-json-without-struct-golang
// We unmarshal yaml into a value of type interface{},
// go through the result recursively, and convert each encountered
// map[interface{}]interface{} to a map[string]interface{} value
// required to marshall to JSON.
// Reference: https://github.com/garethr/kubeval/blob/master/kubeval/utils.go#L8
func convertToStringKeys(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convertToStringKeys(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertToStringKeys(v)
		}
	}
	return i
}

// validate function will validate input kedgefile against JSON schema provided in pkg/validation
func Validate(p []byte) {
	var speco interface{}
	err := yaml.Unmarshal(p, &speco)
	if err != nil {
		fmt.Printf("Error with Unmarhsalling")
	}
	body := convertToStringKeys(speco)
	s := body.(map[string]interface{})
	loader := gojsonschema.NewGoLoader(body)

	if s["controller"] == nil {
		s["controller"] = "deployment"
	}

	var schema gojsonschema.JSONLoader
	// Depend on type of controller, we are fetching jsonschema stored in pkg/validation
	if s["controller"] == "deployment" {
		schema = gojsonschema.NewStringLoader(DeploymentspecmodJson)
	}
	if s["controller"] == "job" {
		schema = gojsonschema.NewStringLoader(JobspecmodJson)
	}
	if s["controller"] == "deploymentconfig" {
		schema = gojsonschema.NewStringLoader(DeploymentconfigspecmodJson)
	}

	// Without forcing these types the schema fails to load
	//Reference: https://github.com/xeipuuv/gojsonschema#formats
	gojsonschema.FormatCheckers.Add("int64", ValidFormat{})
	gojsonschema.FormatCheckers.Add("byte", ValidFormat{})
	gojsonschema.FormatCheckers.Add("int32", ValidFormat{})
	gojsonschema.FormatCheckers.Add("int-or-string", ValidFormat{})
	result, err := gojsonschema.Validate(schema, loader)
	if !result.Valid() {
		fmt.Printf("The kedgefile is not valid. see errors :\n")
		for _, err := range result.Errors() {
			fmt.Printf("- %s\n", err)
			os.Exit(-1)
		}
	}

}
