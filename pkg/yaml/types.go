/*
Copyright 2017 The Kedge Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yaml

import (
	"encoding/json"
	"errors"
	"strconv"
)

// SliceArrayorIntArray represents a slice or an integer
type SliceArrayorIntArray []string

// UnmarshalJSON implements the Unmarshaller interface for the "yaml" package.
// (since the yaml package actually used UnmarshalJSON underneath)
// Overriding the default unmarshaller for our special one.
func (s *SliceArrayorIntArray) UnmarshalJSON(data []byte) error {

	// We create an array of strings by extracting to an interface and then iterating and figuring
	// out each type. As per: https://golang.org/pkg/encoding/json/#Unmarshal
	// we have a "case" for "string" and "float64"
	var unknownType []interface{}
	if err := json.Unmarshal(data, &unknownType); err == nil {

		var stringArray []string

		for _, i := range unknownType {
			switch v := i.(type) {
			case float64:
				// float64, for JSON numbers (int)
				stringArray = append(stringArray, strconv.FormatFloat(v, 'f', 0, 64))
			case string:
				// string, for JSON strings (string)
				stringArray = append(stringArray, v)
			default:
				return errors.New("Unable to convert PortMappings value, are you using a string/int?")
			}
		}

		*s = SliceArrayorIntArray(stringArray)
		return nil
	}

	// If for some reason, we're unable to convert, go to err.
	return errors.New("Failed to unmarshal SliceArrayorIntArray")
}
