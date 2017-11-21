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
	"reflect"
	"testing"
)

type SliceArrayorIntArrayHolder struct {
	SliceorArray SliceArrayorIntArray `json:"portMappings"`
}

func TestSliceArrayorIntArrayUnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		result SliceArrayorIntArray
	}{
		{"{\"portMappings\":[3306]}", SliceArrayorIntArray([]string{"3306"})},
		{"{\"portMappings\":[\"3306\"]}", SliceArrayorIntArray([]string{"3306"})},
		{"{\"portMappings\":[3306, \"3306\"]}", SliceArrayorIntArray([]string{"3306", "3306"})},
		{"{\"portMappings\":[\"80:80\", \"443:443\"]}", SliceArrayorIntArray([]string{"80:80", "443:443"})},
	}

	for _, c := range cases {
		var result SliceArrayorIntArrayHolder
		if err := json.Unmarshal([]byte(c.input), &result); err != nil {
			t.Errorf("Failed to unmarshal input '%v': %v", c.input, err)
		}
		if !reflect.DeepEqual(result.SliceorArray, c.result) {
			t.Errorf("Failed to unmarshal input '%v': expected %+v, got %+v", c.input, c.result, result)
		}
	}
}
