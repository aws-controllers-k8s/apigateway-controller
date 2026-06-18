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

package rest_api

import (
	"testing"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigatewaytypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stretchr/testify/assert"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
)

func TestUpdateRestAPIInput(t *testing.T) {
	for _, tt := range []struct {
		description      string
		desired          svcapitypes.RestAPISpec
		latest           svcapitypes.RestAPISpec
		deltaPaths       []string
		expectedPatchOps []apigatewaytypes.PatchOperation
	}{
		{
			description: "endpointAccessMode change is patched",
			desired:     svcapitypes.RestAPISpec{EndpointAccessMode: aws.String("STRICT")},
			latest:      svcapitypes.RestAPISpec{EndpointAccessMode: aws.String("BASIC")},
			deltaPaths:  []string{"Spec.EndpointAccessMode"},
			expectedPatchOps: []apigatewaytypes.PatchOperation{
				{
					Op:    apigatewaytypes.OpReplace,
					Path:  aws.String("/endpointAccessMode"),
					Value: aws.String("STRICT"),
				},
			},
		},
		{
			description: "securityPolicy change is patched",
			desired:     svcapitypes.RestAPISpec{SecurityPolicy: aws.String("TLS_1_2")},
			latest:      svcapitypes.RestAPISpec{SecurityPolicy: aws.String("TLS_1_0")},
			deltaPaths:  []string{"Spec.SecurityPolicy"},
			expectedPatchOps: []apigatewaytypes.PatchOperation{
				{
					Op:    apigatewaytypes.OpReplace,
					Path:  aws.String("/securityPolicy"),
					Value: aws.String("TLS_1_2"),
				},
			},
		},
		{
			description: "endpointConfiguration ipAddressType change is patched",
			desired: svcapitypes.RestAPISpec{
				EndpointConfiguration: &svcapitypes.EndpointConfiguration{
					IPAddressType: aws.String("dualstack"),
				},
			},
			latest: svcapitypes.RestAPISpec{
				EndpointConfiguration: &svcapitypes.EndpointConfiguration{
					IPAddressType: aws.String("ipv4"),
				},
			},
			deltaPaths: []string{"Spec.EndpointConfiguration.IPAddressType"},
			expectedPatchOps: []apigatewaytypes.PatchOperation{
				{
					Op:    apigatewaytypes.OpReplace,
					Path:  aws.String("/endpointConfiguration/ipAddressType"),
					Value: aws.String("dualstack"),
				},
			},
		},
		{
			description:      "no delta produces no patch operations",
			desired:          svcapitypes.RestAPISpec{SecurityPolicy: aws.String("TLS_1_2")},
			latest:           svcapitypes.RestAPISpec{SecurityPolicy: aws.String("TLS_1_2")},
			deltaPaths:       []string{},
			expectedPatchOps: nil,
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			delta := compare.NewDelta()
			for _, p := range tt.deltaPaths {
				delta.Add(p, nil, nil)
			}
			desired := &resource{ko: &svcapitypes.RestAPI{Spec: tt.desired}}
			latest := &resource{ko: &svcapitypes.RestAPI{Spec: tt.latest}}
			input := &apigateway.UpdateRestApiInput{}

			err := updateRestAPIInput(desired, latest, input, delta)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPatchOps, input.PatchOperations)
		})
	}
}
