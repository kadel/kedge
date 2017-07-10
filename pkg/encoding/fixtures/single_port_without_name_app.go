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

package fixtures

import (
	"github.com/kedgeproject/kedge/pkg/spec"
	api_v1 "k8s.io/client-go/pkg/api/v1"
)

var SinglePortWithoutNameApp spec.App = spec.App{
	Name: "test",
	PodSpecMod: spec.PodSpecMod{
		Containers: []spec.Container{
			{
				Container: api_v1.Container{
					Image: "nginx",
				},
			},
		},
	},
	Services: []spec.ServiceSpecMod{
		{
			Name: "test",
			Ports: []spec.ServicePortMod{
				{
					ServicePort: api_v1.ServicePort{
						Port: 8080,
					},
				},
			},
		},
	},
}
