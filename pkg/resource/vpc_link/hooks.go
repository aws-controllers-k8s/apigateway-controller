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

package vpc_link

import (
	"fmt"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/tags"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

var syncTags = tags.SyncTags

func arnForResource(desired *svcapitypes.VPCLink) (string, error) {
	return util.ARNForResource(desired.Status.ACKResourceMetadata, fmt.Sprintf("/vpclinks/%s", *desired.Status.ID))
}

func validateDeleteState(r *resource) error {
	if status := r.ko.Status.Status; status != nil {
		switch svcapitypes.VPCLinkStatus_SDK(*status) {
		case svcapitypes.VPCLinkStatus_SDK_DELETING, svcapitypes.VPCLinkStatus_SDK_PENDING:
			return ackrequeue.NeededAfter(
				fmt.Errorf("VPCLink is in %s state, it cannot be modified or deleted", *status),
				ackrequeue.DefaultRequeueAfterDuration,
			)
		}
	}
	return nil
}

func updateVPCLinkInput(desired *resource, input *apigateway.UpdateVpcLinkInput, delta *compare.Delta) {
	var patchSet patch.Set
	if delta.DifferentAt("Spec.Name") {
		patchSet.Replace("/name", desired.ko.Spec.Name)
	}
	if delta.DifferentAt("Spec.Description") {
		patchSet.Replace("/description", desired.ko.Spec.Description)
	}
	input.PatchOperations = patchSet.GetPatchOperations()
}
