/*
Copyright 2019 The Crossplane Authors.

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

package resource

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SecretTypeConnection is the type of Crossplane connection secrets.
const SecretTypeConnection corev1.SecretType = "connection.crossplane.io/v1alpha1"

// External resources are tagged/labelled with the following keys in the cloud
// provider API if the type supports.
const (
	ExternalResourceTagKeyKind     = "crossplane-kind"
	ExternalResourceTagKeyName     = "crossplane-name"
	ExternalResourceTagKeyProvider = "crossplane-providerconfig"
)

// A ManagedKind contains the type metadata for a kind of managed resource.
type DnaKind schema.GroupVersionKind

// MustCreateObject returns a new Object of the supplied kind. It panics if the
// kind is unknown to the supplied ObjectCreator.
func MustCreateObject(kind schema.GroupVersionKind, oc runtime.ObjectCreater) runtime.Object {
	obj, err := oc.New(kind)
	if err != nil {
		panic(err)
	}
	return obj
}
