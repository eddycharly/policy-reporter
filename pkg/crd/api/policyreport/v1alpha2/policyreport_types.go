/*
Copyright 2020 The Kubernetes authors.
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

package v1alpha2

import (
	"strconv"

	"github.com/segmentio/fasthash/fnv1a"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kyverno/policy-reporter/pkg/helper"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Kind",type=string,JSONPath=`.scope.kind`,priority=1
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.scope.name`,priority=1
// +kubebuilder:printcolumn:name="Pass",type=integer,JSONPath=`.summary.pass`
// +kubebuilder:printcolumn:name="Fail",type=integer,JSONPath=`.summary.fail`
// +kubebuilder:printcolumn:name="Warn",type=integer,JSONPath=`.summary.warn`
// +kubebuilder:printcolumn:name="Error",type=integer,JSONPath=`.summary.error`
// +kubebuilder:printcolumn:name="Skip",type=integer,JSONPath=`.summary.skip`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName=polr

// PolicyReport is the Schema for the policyreports API
type PolicyReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Scope is an optional reference to the report scope (e.g. a Deployment, Namespace, or Node)
	// +optional
	Scope *corev1.ObjectReference `json:"scope,omitempty"`

	// ScopeSelector is an optional selector for multiple scopes (e.g. Pods).
	// Either one of, or none of, but not both of, Scope or ScopeSelector should be specified.
	// +optional
	ScopeSelector *metav1.LabelSelector `json:"scopeSelector,omitempty"`

	// PolicyReportSummary provides a summary of results
	// +optional
	Summary PolicyReportSummary `json:"summary,omitempty"`

	// PolicyReportResult provides result details
	// +optional
	Results []PolicyReportResult `json:"results,omitempty"`
}

func (r *PolicyReport) GetResults() []PolicyReportResult {
	return r.Results
}

func (r *PolicyReport) SetResults(results []PolicyReportResult) {
	r.Results = results
}

func (r *PolicyReport) GetSummary() PolicyReportSummary {
	return r.Summary
}

func (r *PolicyReport) GetSource() string {
	if len(r.Results) == 0 {
		return ""
	}

	return r.Results[0].Source
}

func (r *PolicyReport) GetKinds() []string {
	list := make([]string, 0)
	for _, k := range r.Results {
		if !k.HasResource() {
			continue
		}

		kind := k.GetResource().Kind

		if kind == "" || helper.Contains(kind, list) {
			continue
		}

		list = append(list, kind)
	}

	return list
}

func (r *PolicyReport) GetSeverities() []string {
	list := make([]string, 0)
	for _, k := range r.Results {

		if k.Severity == "" || helper.Contains(string(k.Severity), list) {
			continue
		}

		list = append(list, string(k.Severity))
	}

	return list
}

func (r *PolicyReport) GetID() string {
	h1 := fnv1a.Init64
	h1 = fnv1a.AddString64(h1, r.GetName())
	h1 = fnv1a.AddString64(h1, r.GetNamespace())

	return strconv.FormatUint(h1, 10)
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyReportList contains a list of PolicyReport
type PolicyReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PolicyReport `json:"items"`
}
