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

package api_method_response

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

var patchKeyEncoder = strings.NewReplacer("~", "~0", "/", "~1")

func updateMethodResponseInput(desired, latest *resource, input *svcsdk.UpdateMethodResponseInput, delta *compare.Delta) {
	latestSpec := latest.ko.Spec
	desiredSpec := desired.ko.Spec

	var patchSet patch.Set
	if delta.DifferentAt("Spec.ResponseParameters") {
		// Handle boolean map patching
		for k := range latestSpec.ResponseParameters {
			if _, ok := desiredSpec.ResponseParameters[k]; !ok {
				patchSet.Remove(fmt.Sprintf("/responseParameters/%s", patchKeyEncoder.Replace(k)), nil)
			}
		}
		for k, v := range desiredSpec.ResponseParameters {
			var strVal *string
			if v != nil {
				strVal = aws.String(strconv.FormatBool(*v))
			}

			if _, ok := latestSpec.ResponseParameters[k]; !ok {
				// Use Add operation for new keys
				patchSet.Add(fmt.Sprintf("/responseParameters/%s", patchKeyEncoder.Replace(k)), strVal)
			} else {
				// Use Replace operation for existing keys
				patchSet.Replace(fmt.Sprintf("/responseParameters/%s", patchKeyEncoder.Replace(k)), strVal)
			}
		}
	}
	if delta.DifferentAt("Spec.ResponseModels") {
		patchSet.ForMap("/responseModels", latestSpec.ResponseModels, desiredSpec.ResponseModels, true)
	}

	input.PatchOperations = patchSet.GetPatchOperations()
}
