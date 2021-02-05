/*


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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaDashboardSpec defines the desired state of GrafanaDashboard
type GrafanaDashboardSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Folder   string         `json:"folder,omitempty"`
	Title    string         `json:"title,omitempty"`
	Editable bool           `json:"editable,omitempty"`
	Rows     []DashboardRow `json:"rows,omitempty"`
}

type DashboardRow struct {
	Name   string           `json:"name,omitempty"`
	Repeat string           `json:"repeat,omitempty"`
	Panels []DashboardPanel `json:"panels,omitempty"`
}

type DashboardPanel struct {
	Title      string             `json:"title,omitempty"`
	Type       string             `json:"type,omitempty"`
	Datasource string             `json:"datasource,omitempty"`
	Targets    []PrometheusTarget `json:"targets,omitempty"`
}

type PrometheusTarget struct {
	Query  string `json:"query,omitempty"`
	Legend string `json:"legend,omitempty"`
	Ref    string `json:"ref,omitempty"`
	Hidden bool   `json:"hidden,omitempty"`
}

// GrafanaDashboardStatus defines the observed state of GrafanaDashboard
type GrafanaDashboardStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status     string `json:"status,omitempty"`
	RetryTimes int    `json:"retryTimes,omitempt"`
}

// +kubebuilder:object:root=true

// GrafanaDashboard is the Schema for the grafanadashboards API
type GrafanaDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaDashboardSpec   `json:"spec,omitempty"`
	Status GrafanaDashboardStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GrafanaDashboardList contains a list of GrafanaDashboard
type GrafanaDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaDashboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrafanaDashboard{}, &GrafanaDashboardList{})
}
