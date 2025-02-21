/*
Copyright 2025.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE! THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required. Any new fields you add
// must have json tags for the fields to be serialized.

// Selector is a simple label selector that matches labels based on a map
// of key-value pairs.
type Selector struct {
	// MatchLabels is a map of {key,value} pairs. A single {key,value}
	// in the matchLabels map is equivalent to an element of matchExpressions,
	// whose key field is "key", the operator is "Equals", and the values array
	// contains only "value". The requirements are ANDed.
	// +kubebuilder:validation:Required
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}

// MachinePoolSpec defines the desired state of a MachinePool.
type MachinePoolSpec struct {
	// Selector is a label query over a set of Machines.
	// The result of matchLabels and matchFields are ANDed.
	// +kubebuilder:validation:Required
	Selector Selector `json:"selector"`
}

// MachinePoolStatus defines the observed state of a MachinePool.
type MachinePoolStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// TODO: Add fields machineCount and readyMachineCount.
	// TODO: Add field for machine names.
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MachinePool is the Schema for the machinepools API
type MachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachinePoolSpec   `json:"spec,omitempty"`
	Status MachinePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachinePoolList contains a list of MachinePool
type MachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachinePool{}, &MachinePoolList{})
}
