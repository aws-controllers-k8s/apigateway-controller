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

// Code generated by ack-generate. DO NOT EDIT.

package rest_api

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	smithy "github.com/aws/smithy-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &svcsdk.Client{}
	_ = &svcapitypes.RestAPI{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
	_ = &aws.Config{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.GetRestApiOutput
	resp, err = rm.sdkapi.GetRestApi(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "GetRestApi", err)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NotFoundException" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.ApiKeySource != "" {
		ko.Spec.APIKeySource = aws.String(string(resp.ApiKeySource))
	} else {
		ko.Spec.APIKeySource = nil
	}
	if resp.BinaryMediaTypes != nil {
		ko.Spec.BinaryMediaTypes = aws.StringSlice(resp.BinaryMediaTypes)
	} else {
		ko.Spec.BinaryMediaTypes = nil
	}
	if resp.CreatedDate != nil {
		ko.Status.CreatedDate = &metav1.Time{*resp.CreatedDate}
	} else {
		ko.Status.CreatedDate = nil
	}
	if resp.Description != nil {
		ko.Spec.Description = resp.Description
	} else {
		ko.Spec.Description = nil
	}
	ko.Spec.DisableExecuteAPIEndpoint = &resp.DisableExecuteApiEndpoint
	if resp.EndpointConfiguration != nil {
		f5 := &svcapitypes.EndpointConfiguration{}
		if resp.EndpointConfiguration.Types != nil {
			f5f0 := []*string{}
			for _, f5f0iter := range resp.EndpointConfiguration.Types {
				var f5f0elem *string
				f5f0elem = aws.String(string(f5f0iter))
				f5f0 = append(f5f0, f5f0elem)
			}
			f5.Types = f5f0
		}
		if resp.EndpointConfiguration.VpcEndpointIds != nil {
			f5.VPCEndpointIDs = aws.StringSlice(resp.EndpointConfiguration.VpcEndpointIds)
		}
		ko.Spec.EndpointConfiguration = f5
	} else {
		ko.Spec.EndpointConfiguration = nil
	}
	if resp.Id != nil {
		ko.Status.ID = resp.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.MinimumCompressionSize != nil {
		minimumCompressionSizeCopy := int64(*resp.MinimumCompressionSize)
		ko.Spec.MinimumCompressionSize = &minimumCompressionSizeCopy
	} else {
		ko.Spec.MinimumCompressionSize = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.Policy != nil {
		ko.Spec.Policy = resp.Policy
	} else {
		ko.Spec.Policy = nil
	}
	if resp.RootResourceId != nil {
		ko.Status.RootResourceID = resp.RootResourceId
	} else {
		ko.Status.RootResourceID = nil
	}
	if resp.Tags != nil {
		ko.Spec.Tags = aws.StringMap(resp.Tags)
	} else {
		ko.Spec.Tags = nil
	}
	if resp.Version != nil {
		ko.Spec.Version = resp.Version
	} else {
		ko.Spec.Version = nil
	}
	if resp.Warnings != nil {
		ko.Status.Warnings = aws.StringSlice(resp.Warnings)
	} else {
		ko.Status.Warnings = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Status.ID == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.GetRestApiInput, error) {
	res := &svcsdk.GetRestApiInput{}

	if r.ko.Status.ID != nil {
		res.RestApiId = r.ko.Status.ID
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateRestApiOutput
	_ = resp
	resp, err = rm.sdkapi.CreateRestApi(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateRestApi", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.ApiKeySource != "" {
		ko.Spec.APIKeySource = aws.String(string(resp.ApiKeySource))
	} else {
		ko.Spec.APIKeySource = nil
	}
	if resp.BinaryMediaTypes != nil {
		ko.Spec.BinaryMediaTypes = aws.StringSlice(resp.BinaryMediaTypes)
	} else {
		ko.Spec.BinaryMediaTypes = nil
	}
	if resp.CreatedDate != nil {
		ko.Status.CreatedDate = &metav1.Time{*resp.CreatedDate}
	} else {
		ko.Status.CreatedDate = nil
	}
	if resp.Description != nil {
		ko.Spec.Description = resp.Description
	} else {
		ko.Spec.Description = nil
	}
	ko.Spec.DisableExecuteAPIEndpoint = &resp.DisableExecuteApiEndpoint
	if resp.EndpointConfiguration != nil {
		f5 := &svcapitypes.EndpointConfiguration{}
		if resp.EndpointConfiguration.Types != nil {
			f5f0 := []*string{}
			for _, f5f0iter := range resp.EndpointConfiguration.Types {
				var f5f0elem *string
				f5f0elem = aws.String(string(f5f0iter))
				f5f0 = append(f5f0, f5f0elem)
			}
			f5.Types = f5f0
		}
		if resp.EndpointConfiguration.VpcEndpointIds != nil {
			f5.VPCEndpointIDs = aws.StringSlice(resp.EndpointConfiguration.VpcEndpointIds)
		}
		ko.Spec.EndpointConfiguration = f5
	} else {
		ko.Spec.EndpointConfiguration = nil
	}
	if resp.Id != nil {
		ko.Status.ID = resp.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.MinimumCompressionSize != nil {
		minimumCompressionSizeCopy := int64(*resp.MinimumCompressionSize)
		ko.Spec.MinimumCompressionSize = &minimumCompressionSizeCopy
	} else {
		ko.Spec.MinimumCompressionSize = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.Policy != nil {
		ko.Spec.Policy = resp.Policy
	} else {
		ko.Spec.Policy = nil
	}
	if resp.RootResourceId != nil {
		ko.Status.RootResourceID = resp.RootResourceId
	} else {
		ko.Status.RootResourceID = nil
	}
	if resp.Tags != nil {
		ko.Spec.Tags = aws.StringMap(resp.Tags)
	} else {
		ko.Spec.Tags = nil
	}
	if resp.Version != nil {
		ko.Spec.Version = resp.Version
	} else {
		ko.Spec.Version = nil
	}
	if resp.Warnings != nil {
		ko.Status.Warnings = aws.StringSlice(resp.Warnings)
	} else {
		ko.Status.Warnings = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateRestApiInput, error) {
	res := &svcsdk.CreateRestApiInput{}

	if r.ko.Spec.APIKeySource != nil {
		res.ApiKeySource = svcsdktypes.ApiKeySourceType(*r.ko.Spec.APIKeySource)
	}
	if r.ko.Spec.BinaryMediaTypes != nil {
		res.BinaryMediaTypes = aws.ToStringSlice(r.ko.Spec.BinaryMediaTypes)
	}
	if r.ko.Spec.CloneFrom != nil {
		res.CloneFrom = r.ko.Spec.CloneFrom
	}
	if r.ko.Spec.Description != nil {
		res.Description = r.ko.Spec.Description
	}
	if r.ko.Spec.DisableExecuteAPIEndpoint != nil {
		res.DisableExecuteApiEndpoint = *r.ko.Spec.DisableExecuteAPIEndpoint
	}
	if r.ko.Spec.EndpointConfiguration != nil {
		f5 := &svcsdktypes.EndpointConfiguration{}
		if r.ko.Spec.EndpointConfiguration.Types != nil {
			f5f0 := []svcsdktypes.EndpointType{}
			for _, f5f0iter := range r.ko.Spec.EndpointConfiguration.Types {
				var f5f0elem string
				f5f0elem = string(*f5f0iter)
				f5f0 = append(f5f0, svcsdktypes.EndpointType(f5f0elem))
			}
			f5.Types = f5f0
		}
		if r.ko.Spec.EndpointConfiguration.VPCEndpointIDs != nil {
			f5.VpcEndpointIds = aws.ToStringSlice(r.ko.Spec.EndpointConfiguration.VPCEndpointIDs)
		}
		res.EndpointConfiguration = f5
	}
	if r.ko.Spec.MinimumCompressionSize != nil {
		minimumCompressionSizeCopy0 := *r.ko.Spec.MinimumCompressionSize
		if minimumCompressionSizeCopy0 > math.MaxInt32 || minimumCompressionSizeCopy0 < math.MinInt32 {
			return nil, fmt.Errorf("error: field MinimumCompressionSize is of type int32")
		}
		minimumCompressionSizeCopy := int32(minimumCompressionSizeCopy0)
		res.MinimumCompressionSize = &minimumCompressionSizeCopy
	}
	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}
	if r.ko.Spec.Policy != nil {
		res.Policy = r.ko.Spec.Policy
	}
	if r.ko.Spec.Tags != nil {
		res.Tags = aws.ToStringMap(r.ko.Spec.Tags)
	}
	if r.ko.Spec.Version != nil {
		res.Version = r.ko.Spec.Version
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()
	if immutableFieldChanges := rm.getImmutableFieldChanges(delta); len(immutableFieldChanges) > 0 {
		msg := fmt.Sprintf("Immutable Spec fields have been modified: %s", strings.Join(immutableFieldChanges, ","))
		return nil, ackerr.NewTerminalError(fmt.Errorf(msg))
	}
	if delta.DifferentAt("Spec.Tags") {
		resourceARN, err := arnForResource(desired.ko)
		if err != nil {
			return nil, fmt.Errorf("applying tags: %w", err)
		}
		if err := syncTags(ctx, rm.sdkapi, rm.metrics, resourceARN, desired.ko.Spec.Tags, latest.ko.Spec.Tags); err != nil {
			return nil, err
		}
	}
	if !delta.DifferentExcept("Spec.Tags") {
		return desired, nil
	}

	input, err := rm.newUpdateRequestPayload(ctx, desired, delta)
	if err != nil {
		return nil, err
	}
	if err := updateRestAPIInput(desired, latest, input, delta); err != nil {
		return nil, err
	}

	var resp *svcsdk.UpdateRestApiOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateRestApi(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateRestApi", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.ApiKeySource != "" {
		ko.Spec.APIKeySource = aws.String(string(resp.ApiKeySource))
	} else {
		ko.Spec.APIKeySource = nil
	}
	if resp.BinaryMediaTypes != nil {
		ko.Spec.BinaryMediaTypes = aws.StringSlice(resp.BinaryMediaTypes)
	} else {
		ko.Spec.BinaryMediaTypes = nil
	}
	if resp.CreatedDate != nil {
		ko.Status.CreatedDate = &metav1.Time{*resp.CreatedDate}
	} else {
		ko.Status.CreatedDate = nil
	}
	if resp.Description != nil {
		ko.Spec.Description = resp.Description
	} else {
		ko.Spec.Description = nil
	}
	ko.Spec.DisableExecuteAPIEndpoint = &resp.DisableExecuteApiEndpoint
	if resp.EndpointConfiguration != nil {
		f5 := &svcapitypes.EndpointConfiguration{}
		if resp.EndpointConfiguration.Types != nil {
			f5f0 := []*string{}
			for _, f5f0iter := range resp.EndpointConfiguration.Types {
				var f5f0elem *string
				f5f0elem = aws.String(string(f5f0iter))
				f5f0 = append(f5f0, f5f0elem)
			}
			f5.Types = f5f0
		}
		if resp.EndpointConfiguration.VpcEndpointIds != nil {
			f5.VPCEndpointIDs = aws.StringSlice(resp.EndpointConfiguration.VpcEndpointIds)
		}
		ko.Spec.EndpointConfiguration = f5
	} else {
		ko.Spec.EndpointConfiguration = nil
	}
	if resp.Id != nil {
		ko.Status.ID = resp.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.MinimumCompressionSize != nil {
		minimumCompressionSizeCopy := int64(*resp.MinimumCompressionSize)
		ko.Spec.MinimumCompressionSize = &minimumCompressionSizeCopy
	} else {
		ko.Spec.MinimumCompressionSize = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.Policy != nil {
		ko.Spec.Policy = resp.Policy
	} else {
		ko.Spec.Policy = nil
	}
	if resp.RootResourceId != nil {
		ko.Status.RootResourceID = resp.RootResourceId
	} else {
		ko.Status.RootResourceID = nil
	}
	if resp.Tags != nil {
		ko.Spec.Tags = aws.StringMap(resp.Tags)
	} else {
		ko.Spec.Tags = nil
	}
	if resp.Version != nil {
		ko.Spec.Version = resp.Version
	} else {
		ko.Spec.Version = nil
	}
	if resp.Warnings != nil {
		ko.Status.Warnings = aws.StringSlice(resp.Warnings)
	} else {
		ko.Status.Warnings = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (*svcsdk.UpdateRestApiInput, error) {
	res := &svcsdk.UpdateRestApiInput{}

	if r.ko.Status.ID != nil {
		res.RestApiId = r.ko.Status.ID
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteRestApiOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteRestApi(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteRestApi", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteRestApiInput, error) {
	res := &svcsdk.DeleteRestApiInput{}

	if r.ko.Status.ID != nil {
		res.RestApiId = r.ko.Status.ID
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.RestAPI,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}

	var terminalErr smithy.APIError
	if !errors.As(err, &terminalErr) {
		return false
	}
	switch terminalErr.ErrorCode() {
	case "BadRequestException",
		"ConflictException",
		"NotFoundException",
		"InvalidParameter":
		return true
	default:
		return false
	}
}

// getImmutableFieldChanges returns list of immutable fields from the
func (rm *resourceManager) getImmutableFieldChanges(
	delta *ackcompare.Delta,
) []string {
	var fields []string
	if delta.DifferentAt("Spec.CloneFrom") {
		fields = append(fields, "CloneFrom")
	}
	if delta.DifferentAt("Spec.Version") {
		fields = append(fields, "Version")
	}

	return fields
}
