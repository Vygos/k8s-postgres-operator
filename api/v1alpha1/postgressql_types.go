/*
Copyright 2024.

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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Status string

const (
	Running Status = "Running"
	Failed  Status = "Failed"
)

// PostgresSQLSpec defines the desired state of PostgresSQL
type PostgresSQLSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//metav1.ObjectMeta

	Selector *metav1.LabelSelector `json:"selector"`
	Env      []v1.EnvVar           `json:"env"`
}

// PostgresSQLStatus defines the observed state of PostgresSQL
type PostgresSQLStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status `json:"status"`
}

func (p *PostgresSQLStatus) Failed() {
	p.Status = Failed
}

func (p *PostgresSQLStatus) Running() {
	p.Status = Running
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PostgresSQL is the Schema for the postgressqls API
type PostgresSQL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresSQLSpec   `json:"spec,omitempty"`
	Status PostgresSQLStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PostgresSQLList contains a list of PostgresSQL
type PostgresSQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresSQL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgresSQL{}, &PostgresSQLList{})
}
