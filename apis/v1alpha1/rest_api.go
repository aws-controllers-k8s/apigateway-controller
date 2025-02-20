// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RestApiSpec defines the desired state of RestApi.
//
// Represents a REST API.
type RestAPISpec struct {

	// The source of the API key for metering requests according to a usage plan.
	// Valid values are: HEADER to read the API key from the X-API-Key header of
	// a request. AUTHORIZER to read the API key from the UsageIdentifierKey from
	// a custom authorizer.
	APIKeySource *string `json:"apiKeySource,omitempty"`
	// The list of binary media types supported by the RestApi. By default, the
	// RestApi supports only UTF-8-encoded text payloads.
	BinaryMediaTypes []*string `json:"binaryMediaTypes,omitempty"`
	// The ID of the RestApi that you want to clone from.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	CloneFrom *string `json:"cloneFrom,omitempty"`
	// The description of the RestApi.
	Description *string `json:"description,omitempty"`
	// Specifies whether clients can invoke your API by using the default execute-api
	// endpoint. By default, clients can invoke your API with the default https://{api_id}.execute-api.{region}.amazonaws.com
	// endpoint. To require that clients use a custom domain name to invoke your
	// API, disable the default endpoint
	DisableExecuteAPIEndpoint *bool `json:"disableExecuteAPIEndpoint,omitempty"`
	// The endpoint configuration of this RestApi showing the endpoint types of
	// the API.
	EndpointConfiguration *EndpointConfiguration `json:"endpointConfiguration,omitempty"`
	// A nullable integer that is used to enable compression (with non-negative
	// between 0 and 10485760 (10M) bytes, inclusive) or disable compression (with
	// a null value) on an API. When compression is enabled, compression or decompression
	// is not applied on the payload if the payload size is smaller than this value.
	// Setting it to zero allows compression for any payload size.
	MinimumCompressionSize *int64 `json:"minimumCompressionSize,omitempty"`
	// The name of the RestApi.
	// +kubebuilder:validation:Required
	Name *string `json:"name"`
	// A stringified JSON policy document that applies to this RestApi regardless
	// of the caller and Method configuration.
	Policy *string `json:"policy,omitempty"`
	// The key-value map of strings. The valid character set is [a-zA-Z+-=._:/].
	// The tag key can be up to 128 characters and must not start with aws:. The
	// tag value can be up to 256 characters.
	Tags map[string]*string `json:"tags,omitempty"`
	// A version identifier for the API.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	Version *string `json:"version,omitempty"`
}

// RestAPIStatus defines the observed state of RestAPI
type RestAPIStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRs managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// The timestamp when the API was created.
	// +kubebuilder:validation:Optional
	CreatedDate *metav1.Time `json:"createdDate,omitempty"`
	// The API's identifier. This identifier is unique across all of your APIs in
	// API Gateway.
	// +kubebuilder:validation:Optional
	ID *string `json:"id,omitempty"`
	// The API's root resource ID.
	// +kubebuilder:validation:Optional
	RootResourceID *string `json:"rootResourceID,omitempty"`
	// The warning messages reported when failonwarnings is turned on during API
	// import.
	// +kubebuilder:validation:Optional
	Warnings []*string `json:"warnings,omitempty"`
}

// RestAPI is the Schema for the RestAPIS API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type RestAPI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RestAPISpec   `json:"spec,omitempty"`
	Status            RestAPIStatus `json:"status,omitempty"`
}

// RestAPIList contains a list of RestAPI
// +kubebuilder:object:root=true
type RestAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RestAPI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RestAPI{}, &RestAPIList{})
}
