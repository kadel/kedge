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

package encoding

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/kedgeproject/kedge/pkg/spec"
)

func Decode(data []byte) (*spec.App, error) {

	var app spec.App
	err := yaml.Unmarshal(data, &app)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmarshalled: %#v\n", app)
	if err := fixApp(&app); err != nil {
		return nil, errors.Wrapf(err, "Unable to fix app %q", app.Name)
	}
	return &app, nil

}
