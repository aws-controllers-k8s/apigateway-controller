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

package method

import (
	"strconv"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

// updateMethodInput adds patch operations to the UpdateMethodInput based on
// differences between the desired and latest states
func updateMethodInput(desired, latest *resource, input *svcsdk.UpdateMethodInput, delta *compare.Delta) {
	latestSpec := latest.ko.Spec
	desiredSpec := desired.ko.Spec

	var patchSet patch.Set

	if delta.DifferentAt("Spec.AuthorizationScopes") {
		patchSet.ForSlice("/authorizationScopes", latestSpec.AuthorizationScopes, desiredSpec.AuthorizationScopes)
	}

	if delta.DifferentAt("Spec.AuthorizationType") {
		patchSet.Replace("/authorizationType", desiredSpec.AuthorizationType)
	}
	if delta.DifferentAt("Spec.AuthorizerID") {
		patchSet.Replace("/authorizerId", desiredSpec.AuthorizerID)
	}
	if delta.DifferentAt("Spec.APIKeyRequired") {
		if desiredSpec.APIKeyRequired != nil {
			val := aws.String(strconv.FormatBool(*desiredSpec.APIKeyRequired))
			patchSet.Replace("/apiKeyRequired", val)
		}
	}
	if delta.DifferentAt("Spec.OperationName") {
		patchSet.Replace("/operationName", desiredSpec.OperationName)
	}

	if delta.DifferentAt("Spec.RequestParameters") {
		latestMap := convertBoolMapToStringMap(latestSpec.RequestParameters)
		desiredMap := convertBoolMapToStringMap(desiredSpec.RequestParameters)
		patchSet.ForMap("/requestParameters", latestMap, desiredMap, true)
	}

	if delta.DifferentAt("Spec.RequestModels") {
		patchSet.ForMap("/requestModels", latestSpec.RequestModels, desiredSpec.RequestModels, true)
	}

	if delta.DifferentAt("Spec.RequestValidatorID") {
		patchSet.Replace("/requestValidatorId", desiredSpec.RequestValidatorID)
	}

	input.PatchOperations = patchSet.GetPatchOperations()
}

func convertBoolMapToStringMap(requestParameters map[string]*bool) map[string]*string {
	requestParametersMap := make(map[string]*string)
	for k, v := range requestParameters {
		if v != nil {
			val := aws.String(strconv.FormatBool(*v))
			requestParametersMap[k] = val
		}
	}
	return requestParametersMap
}
