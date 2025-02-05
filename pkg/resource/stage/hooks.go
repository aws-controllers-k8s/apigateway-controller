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

package stage

import (
	"fmt"
	"strconv"

	"github.com/aws-controllers-k8s/runtime/pkg/compare"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go/aws"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/tags"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util"
	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

var syncTags = tags.SyncTags

func arnForResource(desired *svcapitypes.Stage) (string, error) {
	return util.ARNForResource(desired.Status.ACKResourceMetadata,
		fmt.Sprintf("/restapis/%s/stages/%s", *desired.Spec.RestAPIID, *desired.Spec.StageName))
}

func updateStageInput(desired, latest *resource, input *svcsdk.UpdateStageInput, delta *compare.Delta) {
	latestSpec := latest.ko.Spec
	desiredSpec := desired.ko.Spec

	var patchSet patch.Set
	if delta.DifferentAt("Spec.CacheClusterEnabled") {
		var val *string
		if desiredSpec.CacheClusterEnabled != nil {
			val = aws.String(strconv.FormatBool(*desiredSpec.CacheClusterEnabled))
		}
		patchSet.Replace("/cacheClusterEnabled", val)
	}
	if delta.DifferentAt("Spec.CacheClusterSize") {
		patchSet.Replace("/cacheClusterSize", desiredSpec.CacheClusterSize)
	}
	if delta.DifferentAt("Spec.CanarySettings") {
		updateCanarySettings(delta, desiredSpec, latestSpec, &patchSet)
	}
	if delta.DifferentAt("Spec.DeploymentID") {
		patchSet.Replace("/deploymentId", desiredSpec.DeploymentID)
	}
	if delta.DifferentAt("Spec.Description") {
		patchSet.Replace("/description", desiredSpec.Description)
	}
	if delta.DifferentAt("Spec.DocumentationVersion") {
		patchSet.Replace("/documentationVersion", desiredSpec.DocumentationVersion)
	}
	if delta.DifferentAt("Spec.Variables") {
		patchSet.ForMap("/variables", latestSpec.Variables, desiredSpec.Variables, false)
	}
	if delta.DifferentAt("Spec.TracingEnabled") {
		var val *string
		if desiredSpec.TracingEnabled != nil {
			val = aws.String(strconv.FormatBool(*desiredSpec.TracingEnabled))
		}
		patchSet.Replace("/tracingEnabled", val)
	}
	input.PatchOperations = patchSet.GetPatchOperations()
}

func updateCanarySettings(delta *compare.Delta, desiredSpec, latestSpec svcapitypes.StageSpec, patchSet *patch.Set) {
	const rootKey = "/canarySettings"
	canary := desiredSpec.CanarySettings
	if canary == nil {
		patchSet.Remove(rootKey)
		return
	}

	prefixRootKey := func(key string) string {
		return fmt.Sprintf("%s/%s", rootKey, key)
	}
	if delta.DifferentAt("Spec.CanarySettings.DeploymentID") {
		patchSet.Replace(prefixRootKey("deploymentId"), canary.DeploymentID)
	}
	if delta.DifferentAt("Spec.CanarySettings.PercentTraffic") {
		var val *string
		if canary.PercentTraffic != nil {
			val = aws.String(fmt.Sprintf("%f", *canary.PercentTraffic))
		}
		patchSet.Replace(prefixRootKey("percentTraffic"), val)
	}
	if delta.DifferentAt("Spec.CanarySettings.StageVariableOverrides") {
		desiredValues := canary.StageVariableOverrides
		if desiredValues == nil {
			desiredValues = map[string]*string{}
		}
		var currValues map[string]*string
		if latestSpec.CanarySettings != nil && latestSpec.CanarySettings.StageVariableOverrides != nil {
			currValues = latestSpec.CanarySettings.StageVariableOverrides
		} else {
			currValues = map[string]*string{}
		}
		patchSet.ForMap(prefixRootKey("stageVariableOverrides"), currValues, desiredValues, false)
	}
	if delta.DifferentAt("Spec.CanarySettings.UseStageCache") {
		var val *string
		if canary.UseStageCache != nil {
			val = aws.String(strconv.FormatBool(*canary.UseStageCache))
		}
		patchSet.Replace(prefixRootKey("useStageCache"), val)
	}
}

func customPreCompare(a, b *resource) {
	if a.ko.Spec.Variables == nil && b.ko.Spec.Variables != nil {
		a.ko.Spec.Variables = map[string]*string{}
	} else if a.ko.Spec.Variables != nil && b.ko.Spec.Variables == nil {
		b.ko.Spec.Variables = map[string]*string{}
	}
	if a.ko.Spec.CanarySettings == nil && b.ko.Spec.CanarySettings == nil {
		return
	}
	if a.ko.Spec.CanarySettings == nil {
		a.ko.Spec.CanarySettings = &svcapitypes.CanarySettings{}
	}
	if b.ko.Spec.CanarySettings == nil {
		b.ko.Spec.CanarySettings = &svcapitypes.CanarySettings{}
	}
	if a.ko.Spec.CanarySettings.StageVariableOverrides == nil && b.ko.Spec.CanarySettings.StageVariableOverrides != nil {
		a.ko.Spec.CanarySettings.StageVariableOverrides = map[string]*string{}
	} else if a.ko.Spec.CanarySettings.StageVariableOverrides != nil && b.ko.Spec.CanarySettings.StageVariableOverrides == nil {
		b.ko.Spec.CanarySettings.StageVariableOverrides = map[string]*string{}
	}
}
