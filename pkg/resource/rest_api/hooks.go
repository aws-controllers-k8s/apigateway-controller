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
	"errors"
	"fmt"
	"strconv"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/tags"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

var syncTags = tags.SyncTags

func arnForResource(desired *svcapitypes.RestAPI) (string, error) {
	return util.ARNForResource(desired.Status.ACKResourceMetadata, fmt.Sprintf("/restapis/%s", *desired.Status.ID))
}

func updateRestAPIInput(desired, latest *resource, input *apigateway.UpdateRestApiInput, delta *compare.Delta) error {
	latestSpec := latest.ko.Spec
	desiredSpec := desired.ko.Spec

	var patchSet patch.Set
	if delta.DifferentAt("Spec.APIKeySource") {
		patchSet.Replace("/apiKeySource", desiredSpec.APIKeySource)
	}
	if delta.DifferentAt("Spec.BinaryMediaTypes") {
		patchSet.ForSlice("/binaryMediaTypes", latestSpec.BinaryMediaTypes, desiredSpec.BinaryMediaTypes)
	}
	if delta.DifferentAt("Spec.Description") {
		patchSet.Replace("/description", desiredSpec.Description)
	}
	if delta.DifferentAt("Spec.DisableExecuteAPIEndpoint") {
		var disable bool
		if desiredSpec.DisableExecuteAPIEndpoint != nil {
			disable = *desiredSpec.DisableExecuteAPIEndpoint
		}
		patchSet.Replace("/disableExecuteApiEndpoint", aws.String(strconv.FormatBool(disable)))
	}
	if delta.DifferentAt("Spec.EndpointConfiguration.Types") {
		if desiredSpec.EndpointConfiguration == nil {
			return errors.New("spec.endpointConfiguration.types is required")
		}
		if len(desiredSpec.EndpointConfiguration.Types) != 1 {
			return errors.New("spec.endpointConfiguration.types must contain exactly one element")
		}
		patchSet.Replace("/endpointConfiguration/types/0", desiredSpec.EndpointConfiguration.Types[0])

	}
	if delta.DifferentAt("Spec.EndpointConfiguration.VPCEndpointIDs") {
		var (
			currEndpointIDs    []*string
			desiredEndpointIDs []*string
		)
		if latestSpec.EndpointConfiguration != nil {
			currEndpointIDs = latestSpec.EndpointConfiguration.VPCEndpointIDs
		}
		if desiredSpec.EndpointConfiguration != nil {
			desiredEndpointIDs = desiredSpec.EndpointConfiguration.VPCEndpointIDs
		}
		patchSet.ForSlice("/endpointConfiguration/vpcEndpointIds", currEndpointIDs, desiredEndpointIDs)
	}
	if delta.DifferentAt("Spec.MinimumCompressionSize") {
		var val *string
		if desiredSpec.MinimumCompressionSize != nil {
			val = aws.String(strconv.FormatInt(*desiredSpec.MinimumCompressionSize, 10))
		}
		patchSet.Replace("/minimumCompressionSize", val)
	}
	if delta.DifferentAt("Spec.Name") {
		patchSet.Replace("/name", desiredSpec.Name)
	}
	if delta.DifferentAt("Spec.Policy") {
		patchSet.Replace("/policy", desiredSpec.Policy)
	}
	input.PatchOperations = patchSet.GetPatchOperations()
	return nil
}

func customPreCompare(a, b *resource) {
	if a.ko.Spec.EndpointConfiguration == nil && b.ko.Spec.EndpointConfiguration != nil {
		a.ko.Spec.EndpointConfiguration = &svcapitypes.EndpointConfiguration{}
	} else if a.ko.Spec.EndpointConfiguration != nil && b.ko.Spec.EndpointConfiguration == nil {
		b.ko.Spec.EndpointConfiguration = &svcapitypes.EndpointConfiguration{}
	}
}
