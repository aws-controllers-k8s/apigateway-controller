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

package resource

import (
	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	svcsdk "github.com/aws/aws-sdk-go/service/apigateway"

	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

func updateResourceInput(desired *resource, input *svcsdk.UpdateResourceInput, delta *compare.Delta) {
	desiredSpec := desired.ko.Spec
	var patchSet patch.Set
	if delta.DifferentAt("Spec.ParentID") {
		patchSet.Replace("/parentId", desiredSpec.ParentID)
	}
	if delta.DifferentAt("Spec.PathPart") {
		patchSet.Replace("/pathPart", desiredSpec.PathPart)
	}
	input.PatchOperations = patchSet.GetPatchOperations()
}
