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

package authorizer

import (
	"strconv"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

// updateAuthorizerInput patches the resource based on the delta between desired and latest state.
// It uses the patchSet utility to generate add/remove operations for list fields like providerARNs,
// as indicated by the UpdateAuthorizer API documentation.
func updateAuthorizerInput(
	desired *resource,
	latest *resource,
	input *svcsdk.UpdateAuthorizerInput,
	delta *ackcompare.Delta,
) {
	desiredSpec := desired.ko.Spec
	latestSpec := latest.ko.Spec

	var patchSet patch.Set

	if delta.DifferentAt("Spec.AuthorizerURI") {
		patchSet.Replace("/authorizerUri", desiredSpec.AuthorizerURI)
	}

	if delta.DifferentAt("Spec.AuthorizerCredentials") {
		patchSet.Replace("/authorizerCredentials", desiredSpec.AuthorizerCredentials)
	}

	if delta.DifferentAt("Spec.AuthorizerResultTTLInSeconds") {
		val := aws.String(strconv.FormatInt(*desiredSpec.AuthorizerResultTTLInSeconds, 10))
		patchSet.Replace("/authorizerResultTtlInSeconds", val)
	}

	if delta.DifferentAt("Spec.AuthType") {
		patchSet.Replace("/authType", desiredSpec.AuthType)
	}

	if delta.DifferentAt("Spec.IdentitySource") {
		patchSet.Replace("/identitySource", desiredSpec.IdentitySource)
	}

	if delta.DifferentAt("Spec.IdentityValidationExpression") {
		patchSet.Replace("/identityValidationExpression", desiredSpec.IdentityValidationExpression)
	}

	if delta.DifferentAt("Spec.Name") {
		patchSet.Replace("/name", desiredSpec.Name)
	}

	if delta.DifferentAt("Spec.ProviderARNs") {
		updateProviderARNsPatches(&patchSet, latestSpec.ProviderARNs, desiredSpec.ProviderARNs)
	}

	if delta.DifferentAt("Spec.Type") {
		patchSet.Replace("/type", desiredSpec.Type)
	}

	input.PatchOperations = patchSet.GetPatchOperations()
}

// updateProviderARNsPatches generates patch operations for the providerARNs field
func updateProviderARNsPatches(
	patchSet *patch.Set,
	latestARNs, desiredARNs []*string,
) {
	// Convert []*string to map[string]bool
	latestSet := make(map[string]bool)
	for _, arnPtr := range latestARNs {
		if arnPtr != nil {
			latestSet[*arnPtr] = true
		}
	}

	desiredSet := make(map[string]bool)
	for _, arnPtr := range desiredARNs {
		if arnPtr != nil {
			desiredSet[*arnPtr] = true
		}
	}

	// ARNs to remove
	for arn := range latestSet {
		if !desiredSet[arn] {
			// Use RemoveWithValue to generate: op=remove, path=/providerARNs, value=arn
			patchSet.Remove("/providerARNs", aws.String(arn))
		}
	}

	// ARNs to add
	for arn := range desiredSet {
		if !latestSet[arn] {
			patchSet.Add("/providerARNs", aws.String(arn))
		}
	}
}
