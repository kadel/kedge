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

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	//kapi "k8s.io/kubernetes/pkg/api/v1"
	batch_v1 "k8s.io/kubernetes/pkg/apis/batch/v1"
)

// This function will search in the pod level volumes
// and see if the volume with given name is defined
func isVolumeDefined(volumes []api_v1.Volume, name string) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// search through all the persistent volumes defined in the root level
func isPVCDefined(volumes []VolumeClaim, name string) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// GetScheme() returns runtime.Scheme with supported Kubernetes API resource
// definitions which Kedge supports right now.
// The core v1 scheme is first initialized and then other controllers' scheme
// is added to that scheme, e.g. batch/v1 scheme is added to add support for
// Jobs controller to the v1 Scheme.
// Also, (from upstream) Scheme defines methods for serializing and deserializing API objects, a type
// registry for converting group, version, and kind information to and from Go
// schemas, and mappings between Go schemas of different versions. A scheme is the
// foundation for a versioned API and versioned configuration over time.
func GetScheme() (*runtime.Scheme, error) {
	// Initializing the scheme with the core v1 api
	scheme := api.Scheme

	// Adding the batch scheme to support Jobs
	// TODO: find a way where we don't have to add batch/v1 to the v1 scheme,
	// instead we should be able to have different scheme for different controllers
	if err := batch_v1.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add 'batch' to scheme")
	}
	return scheme, nil
}

// SetGVK() sets Group, Version and Kind for the generated Kubernetes resources.
// This takes in a generated Kubernetes API resource's runtime object and
// runtime scheme based on which the GVK will be set.
func SetGVK(runtimeObject runtime.Object, scheme *runtime.Scheme) error {
	gvk, isUnversioned, err := scheme.ObjectKind(runtimeObject)
	if err != nil {
		return errors.Wrap(err, "ConvertToVersion failed")
	}
	if isUnversioned {
		return fmt.Errorf("ConvertToVersion failed: can't output unversioned type: %T", runtimeObject)
	}
	runtimeObject.GetObjectKind().SetGroupVersionKind(gvk)
	return nil
}

func getInt32Addr(i int32) *int32 {
	return &i
}

func getInt64Addr(i int64) *int64 {
	return &i
}

func prettyPrintObjects(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// addKeyValueToMap adds a key value pair to a given map[string]string only if
// the map does not contain the supplied key. Creates a new map if map is empty.
// We need to return the map because in case a nil map is passed to this
// function, the new map created will not be reflected in the original nil map.
func addKeyValueToMap(k string, v string, m map[string]string) map[string]string {

	if len(m) == 0 {
		m = make(map[string]string)
	}

	if _, ok := m[k]; !ok {
		m[k] = v
	} else {
		log.Debugf("not adding '%v: %v' to map since there exists a user defined label '%v: %v'", k, v, k, m[k])
	}

	return m
}

// // Converts from k8s.io/kubernetes/pkg/api/v1 to kubernetes
// // due to kedge requiring kubernetes but OpenShift using k8s/io/kubernetes
// // This function converts a pod spec to kapi.PodSpec
// func ConvertPodSpec(pod api_v1.PodSpec) kapi.PodSpec {

// 	containers := []kapi.Container{}
// 	for _, container := range pod.Containers {

// 		// Add all keys which don't use custom structs
// 		con := kapi.Container{
// 			Name:                     container.Name,
// 			Image:                    container.Image,
// 			Command:                  container.Command,
// 			Args:                     container.Args,
// 			WorkingDir:               container.WorkingDir,
// 			TerminationMessagePath:   container.TerminationMessagePath,
// 			TerminationMessagePolicy: kapi.TerminationMessagePolicy(container.TerminationMessagePolicy),
// 			ImagePullPolicy:          kapi.PullPolicy(container.ImagePullPolicy),
// 			Stdin:                    container.Stdin,
// 			StdinOnce:                container.StdinOnce,
// 			TTY:                      container.TTY,
// 		}

// 		// TODO: (too difficult to implement at the moment, will need to be implemented in the future)
// 		// Resources
// 		// Lifecycle
// 		// SecurityContext

// 		// Ports
// 		for _, port := range container.Ports {
// 			con.Ports = append(con.Ports, kapi.ContainerPort{
// 				Name:          port.Name,
// 				HostPort:      port.HostPort,
// 				ContainerPort: port.ContainerPort,
// 				HostIP:        port.HostIP,
// 				Protocol:      kapi.Protocol(port.Protocol),
// 			})
// 		}

// 		// EnvFrom
// 		for _, envFrom := range container.EnvFrom {

// 			e := kapi.EnvFromSource{
// 				Prefix: envFrom.Prefix,
// 			}

// 			if envFrom.ConfigMapRef != nil {
// 				e.ConfigMapRef = &kapi.ConfigMapEnvSource{
// 					LocalObjectReference: kapi.LocalObjectReference{
// 						Name: envFrom.ConfigMapRef.LocalObjectReference.Name,
// 					},
// 					Optional: envFrom.ConfigMapRef.Optional,
// 				}
// 			}

// 			con.EnvFrom = append(con.EnvFrom, e)
// 		}

// 		// Env
// 		for _, env := range container.Env {
// 			e := kapi.EnvVar{
// 				Name:      env.Name,
// 				Value:     env.Value,
// 				ValueFrom: &kapi.EnvVarSource{},
// 			}

// 			if env.ValueFrom != nil {

// 				if env.ValueFrom.FieldRef != nil {
// 					e.ValueFrom.FieldRef = &kapi.ObjectFieldSelector{
// 						APIVersion: env.ValueFrom.FieldRef.APIVersion,
// 						FieldPath:  env.ValueFrom.FieldRef.FieldPath,
// 					}
// 				}

// 				if env.ValueFrom.ResourceFieldRef != nil {
// 					e.ValueFrom.ResourceFieldRef = &kapi.ResourceFieldSelector{
// 						ContainerName: env.ValueFrom.ResourceFieldRef.ContainerName,
// 						Resource:      env.ValueFrom.ResourceFieldRef.Resource,
// 						// TODO Divisor
// 					}
// 				}

// 				if env.ValueFrom.ConfigMapKeyRef != nil {
// 					e.ValueFrom.ConfigMapKeyRef = &kapi.ConfigMapKeySelector{
// 						LocalObjectReference: kapi.LocalObjectReference{
// 							Name: env.ValueFrom.ConfigMapKeyRef.LocalObjectReference.Name,
// 						},
// 						Key:      env.ValueFrom.ConfigMapKeyRef.Key,
// 						Optional: env.ValueFrom.ConfigMapKeyRef.Optional,
// 					}
// 				}

// 				if env.ValueFrom.SecretKeyRef != nil {
// 					e.ValueFrom.SecretKeyRef = &kapi.SecretKeySelector{
// 						LocalObjectReference: kapi.LocalObjectReference{
// 							Name: env.ValueFrom.SecretKeyRef.LocalObjectReference.Name,
// 						},
// 						Key:      env.ValueFrom.SecretKeyRef.Key,
// 						Optional: env.ValueFrom.SecretKeyRef.Optional,
// 					}
// 				}
// 			}

// 			con.Env = append(con.Env, e)
// 		}

// 		// VolumeMounts
// 		for _, volume := range container.VolumeMounts {
// 			con.VolumeMounts = append(con.VolumeMounts, kapi.VolumeMount{
// 				Name:      volume.Name,
// 				ReadOnly:  volume.ReadOnly,
// 				MountPath: volume.MountPath,
// 				SubPath:   volume.SubPath,
// 			})
// 		}

// 		// LivenessProbe
// 		if container.LivenessProbe != nil {

// 			con.LivenessProbe = &kapi.Probe{
// 				InitialDelaySeconds: container.LivenessProbe.InitialDelaySeconds,
// 				TimeoutSeconds:      container.LivenessProbe.TimeoutSeconds,
// 				PeriodSeconds:       container.LivenessProbe.PeriodSeconds,
// 				SuccessThreshold:    container.LivenessProbe.SuccessThreshold,
// 				FailureThreshold:    container.LivenessProbe.FailureThreshold,
// 			}

// 			con.LivenessProbe.Handler = convertHandler(container.LivenessProbe.Handler)

// 		}

// 		// ReadinessProbe
// 		if container.ReadinessProbe != nil {

// 			con.ReadinessProbe = &kapi.Probe{
// 				InitialDelaySeconds: container.ReadinessProbe.InitialDelaySeconds,
// 				TimeoutSeconds:      container.ReadinessProbe.TimeoutSeconds,
// 				PeriodSeconds:       container.ReadinessProbe.PeriodSeconds,
// 				SuccessThreshold:    container.ReadinessProbe.SuccessThreshold,
// 				FailureThreshold:    container.ReadinessProbe.FailureThreshold,
// 			}

// 			con.ReadinessProbe.Handler = convertHandler(container.ReadinessProbe.Handler)
// 		}

// 		// Add container
// 		containers = append(containers, con)
// 	}

// 	return kapi.PodSpec{
// 		Containers: containers,
// 	}
// }

// // Converts a kubernetes Handler struct to a kubernetes/kubernetes struct
// func convertHandler(handler api_v1.Handler) kapi.Handler {

// 	convertedHandler := kapi.Handler{}

// 	if handler.Exec != nil {
// 		convertedHandler.Exec = &kapi.ExecAction{Command: handler.Exec.Command}
// 	}

// 	if handler.HTTPGet != nil {
// 		convertedHandler.HTTPGet = &kapi.HTTPGetAction{
// 			Path:   handler.HTTPGet.Path,
// 			Port:   intstr.IntOrString(handler.HTTPGet.Port),
// 			Host:   handler.HTTPGet.Host,
// 			Scheme: kapi.URIScheme(handler.HTTPGet.Scheme),
// 		}
// 		for _, header := range handler.HTTPGet.HTTPHeaders {
// 			convertedHandler.HTTPGet.HTTPHeaders = append(convertedHandler.HTTPGet.HTTPHeaders, kapi.HTTPHeader{
// 				Name:  header.Name,
// 				Value: header.Value,
// 			})
// 		}
// 	}

// 	if handler.TCPSocket != nil {
// 		convertedHandler.TCPSocket = &kapi.TCPSocketAction{Port: intstr.IntOrString(handler.TCPSocket.Port)}
// 	}

// 	return convertedHandler
// }
