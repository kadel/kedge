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

package spec

import (
	"fmt"
	"reflect"

	log "github.com/Sirupsen/logrus"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	ext_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
)

func (app *App) validateDeployments() error {

	// TODO: v2
	return nil
}

func (app *App) fixDeployments() error {
	for _, deployment := range app.Deployments {

		deployment.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, app.Name, deployment.ObjectMeta.Labels)

		if app.Appversion != "" {
			deployment.ObjectMeta.Annotations = addKeyValueToMap(appVersion, app.Appversion, deployment.ObjectMeta.Annotations)
		}
	}

	return nil
}

// Creates a Deployment Kubernetes resource. The returned Deployment resource
// will be nil if it could not be generated due to insufficient input data.
func (app *App) createDeployments() ([]runtime.Object, error) {
	var deployments []runtime.Object

	for _, deployment := range app.Deployments {

		// We need to error out if both, deployment.PodSpec and deployment.DeploymentSpec are empty

		if deployment.isDeploymentSpecPodSpecEmpty() {
			log.Debug("Both, deployment.PodSpec and deployment.DeploymentSpec are empty, not enough data to create a deployment.")
			return nil, nil
		}

		// We are merging whole DeploymentSpec with PodSpec.
		// This means that someone could specify containers in template.spec and also in top level PodSpec.
		// This stupid check is supposed to make sure that only one of them set.
		// TODO: merge DeploymentSpec.Template.Spec and top level PodSpec
		if deployment.isMultiplePodSpecSpecified() {
			return nil, fmt.Errorf("Pod can't be specfied in two places. Use top level PodSpec or template.spec (DeploymentSpec.Template.Spec) not both")
		}

		deploymentSpec := deployment.DeploymentSpec

		// top level PodSpec is not empty, use it for deployment template
		// we already know that if deployment.PodSpec is not empty deployment.DeploymentSpec.Template.Spec is empty
		if !reflect.DeepEqual(deployment.PodSpec, api_v1.PodSpec{}) {
			deploymentSpec.Template.Spec = deployment.PodSpec
		}

		// TODO: check if this wasn't set by user, in that case we shouldn't overwrite it
		deploymentSpec.Template.ObjectMeta.Name = deployment.Name

		// TODO: merge with already existing labels and avoid duplication
		deploymentSpec.Template.ObjectMeta.Labels = deployment.Labels

		deploymentSpec.Template.ObjectMeta.Annotations = deployment.Annotations

		deployments = append(deployments, &ext_v1beta1.Deployment{
			ObjectMeta: deployment.ObjectMeta,
			Spec:       deploymentSpec,
		})
	}
	return deployments, nil
}

func (deployment *DeploymentSpecMod) isDeploymentSpecPodSpecEmpty() bool {
	return reflect.DeepEqual(deployment.PodSpec, api_v1.PodSpec{}) && reflect.DeepEqual(deployment.DeploymentSpec, ext_v1beta1.DeploymentSpec{})
}

func (deployment *DeploymentSpecMod) isMultiplePodSpecSpecified() bool {
	return !(reflect.DeepEqual(deployment.DeploymentSpec.Template.Spec, api_v1.PodSpec{}) || reflect.DeepEqual(deployment.PodSpec, api_v1.PodSpec{}))
}
