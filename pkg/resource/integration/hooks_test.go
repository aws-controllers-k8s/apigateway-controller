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

package integration

import (
	"testing"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigatewaytypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stretchr/testify/assert"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
)

func TestUpdateIntegrationInput(t *testing.T) {
	for _, tt := range []struct {
		description      string
		desired          svcapitypes.IntegrationSpec
		latest           svcapitypes.IntegrationSpec
		deltaPaths       []string
		expectedPatchOps []apigatewaytypes.PatchOperation
	}{
		{
			description: "integrationTarget change is patched",
			desired:     svcapitypes.IntegrationSpec{IntegrationTarget: aws.String("arn:listener-new")},
			latest:      svcapitypes.IntegrationSpec{IntegrationTarget: aws.String("arn:listener-old")},
			deltaPaths:  []string{"Spec.IntegrationTarget"},
			expectedPatchOps: []apigatewaytypes.PatchOperation{
				{
					Op:    apigatewaytypes.OpReplace,
					Path:  aws.String("/integrationTarget"),
					Value: aws.String("arn:listener-new"),
				},
			},
		},
		{
			description: "responseTransferMode change is patched",
			desired:     svcapitypes.IntegrationSpec{ResponseTransferMode: aws.String("STREAM")},
			latest:      svcapitypes.IntegrationSpec{ResponseTransferMode: aws.String("BUFFERED")},
			deltaPaths:  []string{"Spec.ResponseTransferMode"},
			expectedPatchOps: []apigatewaytypes.PatchOperation{
				{
					Op:    apigatewaytypes.OpReplace,
					Path:  aws.String("/responseTransferMode"),
					Value: aws.String("STREAM"),
				},
			},
		},
		{
			description:      "no delta produces no patch operations",
			desired:          svcapitypes.IntegrationSpec{IntegrationTarget: aws.String("arn:listener")},
			latest:           svcapitypes.IntegrationSpec{IntegrationTarget: aws.String("arn:listener")},
			deltaPaths:       []string{},
			expectedPatchOps: nil,
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			delta := compare.NewDelta()
			for _, p := range tt.deltaPaths {
				delta.Add(p, nil, nil)
			}
			desired := &resource{ko: &svcapitypes.Integration{Spec: tt.desired}}
			latest := &resource{ko: &svcapitypes.Integration{Spec: tt.latest}}
			input := &svcsdk.UpdateIntegrationInput{}

			updateIntegrationInput(desired, latest, input, delta)

			assert.Equal(t, tt.expectedPatchOps, input.PatchOperations)
		})
	}
}
