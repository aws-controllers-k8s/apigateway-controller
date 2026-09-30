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
	"net/url"
	"regexp"
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
	if delta.DifferentAt("Spec.EndpointAccessMode") {
		patchSet.Replace("/endpointAccessMode", desiredSpec.EndpointAccessMode)
	}
	if delta.DifferentAt("Spec.EndpointConfiguration.IPAddressType") {
		var ipAddressType *string
		if desiredSpec.EndpointConfiguration != nil {
			ipAddressType = desiredSpec.EndpointConfiguration.IPAddressType
		}
		patchSet.Replace("/endpointConfiguration/ipAddressType", ipAddressType)
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
	if delta.DifferentAt("Spec.SecurityPolicy") {
		patchSet.Replace("/securityPolicy", desiredSpec.SecurityPolicy)
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

	// Normalise the policy returned by AWS before the delta comparison.
	// AWS returns the resource policy URL-encoded, and expands the portable
	// shorthand "execute-api:/*" to the concrete RestAPI ARN
	// "arn:aws:execute-api:<region>:<account>:<api-id>/*". Without normalisation
	// these are never equal: the delta fires on every reconcile, creating an
	// infinite patch loop. We normalise the observed (latest) policy only; the
	// desired policy is left as the user authored it.
	if b.ko.Spec.Policy != nil {
		b.ko.Spec.Policy = normalizeRestAPIPolicy(b.ko.Status.ID, b.ko.Spec.Policy)
	}
}

// normalizeRestAPIPolicy URL-decodes the policy string returned by AWS (the
// API Gateway REST API returns it percent-encoded) and replaces any concrete
// execute-api ARN that refers to the current RestAPI with the portable
// shorthand "execute-api:/<resource>", so that the delta check treats the
// shorthand and the AWS-expanded form as equivalent.
func normalizeRestAPIPolicy(apiID *string, policy *string) *string {
	if policy == nil {
		return nil
	}
	decoded, err := url.QueryUnescape(*policy)
	if err != nil {
		decoded = *policy
	}
	if apiID != nil && *apiID != "" {
		re := regexp.MustCompile(
			`arn:aws:execute-api:[^:]+:[^:]+:` + regexp.QuoteMeta(*apiID) + `/`,
		)
		decoded = re.ReplaceAllString(decoded, "execute-api:/")
	}
	return &decoded
}
