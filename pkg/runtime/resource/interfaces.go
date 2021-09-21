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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// An Object is a Kubernetes object.
type Object interface {
	metav1.Object
	runtime.Object
}

// A Dna is a Kubernetes object representing a concrete dna
// resource (e.g. a Firewall rule).
type Dna interface {
	Object
}

// A DnaList is a list of dna resources.
type DnaList interface {
	client.ObjectList

	// GetItems returns the list of dna resources.
	GetItems() []Dna
}
