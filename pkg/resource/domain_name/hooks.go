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

package domain_name

import (
	"context"
	"fmt"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
	svcapitags "github.com/aws-controllers-k8s/apigateway-controller/pkg/tags"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

func updateDomainNameInput(desired, latest *resource, input *svcsdk.UpdateDomainNameInput, delta *compare.Delta) error {
	latestSpec := latest.ko.Spec
	desiredSpec := desired.ko.Spec

	var patchSet patch.Set

	handleEndpointConfigurationChanges(&patchSet, &latestSpec, &desiredSpec, delta)
	handleRegionalCertificates(&patchSet, &latestSpec, &desiredSpec, delta)
	handleOwnershipVerificationCertificate(&patchSet, &latestSpec, &desiredSpec, delta)
	handlePolicyChanges(&patchSet, &desiredSpec, delta)

	input.PatchOperations = patchSet.GetPatchOperations()
	return nil
}

func handleEndpointConfigurationChanges(patchSet *patch.Set, latestSpec, desiredSpec *svcapitypes.DomainNameSpec, delta *compare.Delta) {
	handleStandardCertificateChanges(patchSet, desiredSpec, delta)

	if delta.DifferentAt("Spec.EndpointConfiguration") && delta.DifferentAt("Spec.EndpointConfiguration.Types") {
		var currTypes []*string
		if latestSpec.EndpointConfiguration != nil {
			currTypes = latestSpec.EndpointConfiguration.Types
		}

		var desiredTypes []*string
		if desiredSpec.EndpointConfiguration != nil {
			desiredTypes = desiredSpec.EndpointConfiguration.Types
		}

		patchSet.ForSlice("/endpointConfiguration/types", currTypes, desiredTypes)
	}

	handleVPCEndpointIDs(patchSet, latestSpec, desiredSpec)
}

func handleStandardCertificateChanges(patchSet *patch.Set, desiredSpec *svcapitypes.DomainNameSpec, delta *compare.Delta) {
	if delta.DifferentAt("Spec.CertificateARN") {
		patchSet.Replace("/certificateArn", desiredSpec.CertificateARN)
	}
	if delta.DifferentAt("Spec.CertificateName") {
		patchSet.Replace("/certificateName", desiredSpec.CertificateName)
	}
}

func handleVPCEndpointIDs(patchSet *patch.Set, latestSpec, desiredSpec *svcapitypes.DomainNameSpec) {
	if desiredSpec.EndpointConfiguration.VPCEndpointIDs == nil {
		return
	}

	var currEndpointIDs []*string
	if latestSpec.EndpointConfiguration != nil {
		currEndpointIDs = latestSpec.EndpointConfiguration.VPCEndpointIDs
	}

	patchSet.ForSlice("/endpointConfiguration/vpcEndpointIds",
		currEndpointIDs,
		desiredSpec.EndpointConfiguration.VPCEndpointIDs)
}

func handleRegionalCertificates(patchSet *patch.Set, latestSpec, desiredSpec *svcapitypes.DomainNameSpec, delta *compare.Delta) {
	handleFieldWithAddReplaceRemove(patchSet,
		"/regionalCertificateArn",
		"Spec.RegionalCertificateARN",
		desiredSpec.RegionalCertificateARN,
		latestSpec.RegionalCertificateARN,
		delta)

	handleFieldWithAddReplaceRemove(patchSet,
		"/regionalCertificateName",
		"Spec.RegionalCertificateName",
		desiredSpec.RegionalCertificateName,
		latestSpec.RegionalCertificateName,
		delta)
}

func handleOwnershipVerificationCertificate(patchSet *patch.Set, latestSpec, desiredSpec *svcapitypes.DomainNameSpec, delta *compare.Delta) {
	handleFieldWithAddReplaceRemove(patchSet,
		"/ownershipVerificationCertificateArn",
		"Spec.OwnershipVerificationCertificateARN",
		desiredSpec.OwnershipVerificationCertificateARN,
		latestSpec.OwnershipVerificationCertificateARN,
		delta)
}

func handlePolicyChanges(patchSet *patch.Set, desiredSpec *svcapitypes.DomainNameSpec, delta *compare.Delta) {
	if delta.DifferentAt("Spec.Policy") {
		patchSet.Replace("/policy", desiredSpec.Policy)
	}
	if delta.DifferentAt("Spec.SecurityPolicy") {
		patchSet.Replace("/securityPolicy", desiredSpec.SecurityPolicy)
	}
}

func handleFieldWithAddReplaceRemove(patchSet *patch.Set, path, deltaPath string, desiredValue, latestValue *string, delta *compare.Delta) {
	if !delta.DifferentAt(deltaPath) {
		return
	}
	switch {
	case desiredValue == nil:
		patchSet.Remove(path, nil)
	case latestValue == nil:
		patchSet.Add(path, desiredValue)
	default:
		patchSet.Replace(path, desiredValue)
	}
}

// getDomainNameARN returns the ARN for a given domain name
func (rm *resourceManager) getDomainNameARN(domainName string) string {
	// API Gateway domain name ARN format:
	// arn:aws:apigateway:region::/domainnames/domain-name
	return fmt.Sprintf(
		"arn:%s:apigateway:%s::/domainnames/%s",
		rm.awsAccountID,
		rm.awsRegion,
		domainName,
	)
}

// syncTags synchronizes the tags for a given domain name
func (rm *resourceManager) syncTags(
	ctx context.Context,
	latest *resource,
	desired *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func() { exit(err) }()

	// Get the ARN of the domain name
	resourceARN := rm.getDomainNameARN(*desired.ko.Spec.DomainName)

	// Get the existing tags
	existingTags := map[string]*string{}
	if latest != nil && latest.ko.Spec.Tags != nil {
		existingTags = latest.ko.Spec.Tags
	}

	// Get the desired tags
	desiredTags := map[string]*string{}
	if desired.ko.Spec.Tags != nil {
		desiredTags = desired.ko.Spec.Tags
	}

	return svcapitags.SyncTags(
		ctx,
		rm.sdkapi,
		rm.metrics,
		resourceARN,
		desiredTags,
		existingTags,
	)
}
