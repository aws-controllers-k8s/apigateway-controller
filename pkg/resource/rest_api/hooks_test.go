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
	"net/url"
	"testing"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigatewaytypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stretchr/testify/assert"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
)

const (
	testAPIID = "abc123def0"
	testRegion  = "us-east-1"
	testAccount = "123456789012"
)

// shorthandPolicy returns a minimal resource policy using the execute-api shorthand.
func shorthandPolicy(resource string) string {
	return `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"execute-api:Invoke","Resource":"execute-api:` + resource + `"}]}`
}

// expandedARNPolicy returns the same policy with the concrete RestAPI ARN.
func expandedARNPolicy(resource string) string {
	return `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"execute-api:Invoke","Resource":"arn:aws:execute-api:` + testRegion + `:` + testAccount + `:` + testAPIID + `/` + resource + `"}]}`
}

func TestNormalizeRestAPIPolicy(t *testing.T) {
	for _, tt := range []struct {
		description string
		apiID       *string
		policy      *string
		want        string
	}{
		{
			description: "nil policy returns nil",
			apiID:       aws.String(testAPIID),
			policy:      nil,
		},
		{
			description: "plain shorthand is returned unchanged",
			apiID:       aws.String(testAPIID),
			policy:      aws.String(shorthandPolicy("/*")),
			want:        shorthandPolicy("/*"),
		},
		{
			description: "expanded ARN is normalised to shorthand",
			apiID:       aws.String(testAPIID),
			policy:      aws.String(expandedARNPolicy("*")),
			want:        shorthandPolicy("/*"),
		},
		{
			description: "URL-encoded expanded ARN is decoded and normalised",
			apiID:       aws.String(testAPIID),
			policy:      aws.String(url.QueryEscape(expandedARNPolicy("*"))),
			want:        shorthandPolicy("/*"),
		},
		{
			description: "ARN for a different API ID is not normalised",
			apiID:       aws.String(testAPIID),
			policy:      aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Resource":"arn:aws:execute-api:us-east-1:123456789012:OTHER_API/*"}]}`),
			want:        `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Resource":"arn:aws:execute-api:us-east-1:123456789012:OTHER_API/*"}]}`,
		},
		{
			description: "nil API ID leaves expanded ARN unchanged",
			apiID:       nil,
			policy:      aws.String(expandedARNPolicy("*")),
			want:        expandedARNPolicy("*"),
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			got := normalizeRestAPIPolicy(tt.apiID, tt.policy)
			if tt.policy == nil {
				assert.Nil(t, got)
				return
			}
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestCustomPreCompare_PolicyNormalization(t *testing.T) {
	for _, tt := range []struct {
		description     string
		desiredPolicy   *string
		latestPolicy    *string
		latestAPIID     *string
		expectNoDelta   bool
	}{
		{
			description:   "shorthand vs expanded ARN produces no delta",
			desiredPolicy: aws.String(shorthandPolicy("/*")),
			latestPolicy:  aws.String(expandedARNPolicy("*")),
			latestAPIID:   aws.String(testAPIID),
			expectNoDelta: true,
		},
		{
			description:   "URL-encoded expanded ARN vs shorthand produces no delta",
			desiredPolicy: aws.String(shorthandPolicy("/*")),
			latestPolicy:  aws.String(url.QueryEscape(expandedARNPolicy("*"))),
			latestAPIID:   aws.String(testAPIID),
			expectNoDelta: true,
		},
		{
			description:   "genuine policy drift is still detected",
			desiredPolicy: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"execute-api:Invoke","Resource":"execute-api:/*","Condition":{"IpAddress":{"aws:SourceIp":"10.0.0.0/8"}}}]}`),
			latestPolicy:  aws.String(expandedARNPolicy("*")),
			latestAPIID:   aws.String(testAPIID),
			expectNoDelta: false,
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			desired := &resource{ko: &svcapitypes.RestAPI{
				Spec: svcapitypes.RestAPISpec{Policy: tt.desiredPolicy},
			}}
			latest := &resource{ko: &svcapitypes.RestAPI{
				Spec:   svcapitypes.RestAPISpec{Policy: tt.latestPolicy},
				Status: svcapitypes.RestAPIStatus{ID: tt.latestAPIID},
			}}
			delta := newResourceDelta(desired, latest)
			if tt.expectNoDelta {
				assert.False(t, delta.DifferentAt("Spec.Policy"),
					"expected no Spec.Policy delta but one was found")
			} else {
				assert.True(t, delta.DifferentAt("Spec.Policy"),
					"expected Spec.Policy delta but none was found")
			}
		})
	}
}

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
