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

package authorizer

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
	_ = &svcapitypes.Authorizer{}
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

	var resp *svcsdk.GetAuthorizerOutput
	resp, err = rm.sdkapi.GetAuthorizer(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "GetAuthorizer", err)
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

	if resp.AuthType != nil {
		ko.Spec.AuthType = resp.AuthType
	} else {
		ko.Spec.AuthType = nil
	}
	if resp.AuthorizerCredentials != nil {
		ko.Spec.AuthorizerCredentials = resp.AuthorizerCredentials
	} else {
		ko.Spec.AuthorizerCredentials = nil
	}
	if resp.AuthorizerResultTtlInSeconds != nil {
		authorizerResultTTLInSecondsCopy := int64(*resp.AuthorizerResultTtlInSeconds)
		ko.Spec.AuthorizerResultTTLInSeconds = &authorizerResultTTLInSecondsCopy
	} else {
		ko.Spec.AuthorizerResultTTLInSeconds = nil
	}
	if resp.AuthorizerUri != nil {
		ko.Spec.AuthorizerURI = resp.AuthorizerUri
	} else {
		ko.Spec.AuthorizerURI = nil
	}
	if resp.Id != nil {
		ko.Status.ID = resp.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.IdentitySource != nil {
		ko.Spec.IdentitySource = resp.IdentitySource
	} else {
		ko.Spec.IdentitySource = nil
	}
	if resp.IdentityValidationExpression != nil {
		ko.Spec.IdentityValidationExpression = resp.IdentityValidationExpression
	} else {
		ko.Spec.IdentityValidationExpression = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ProviderARNs != nil {
		ko.Spec.ProviderARNs = aws.StringSlice(resp.ProviderARNs)
	} else {
		ko.Spec.ProviderARNs = nil
	}
	if resp.Type != "" {
		ko.Spec.Type = aws.String(string(resp.Type))
	} else {
		ko.Spec.Type = nil
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
	return r.ko.Status.ID == nil || r.ko.Spec.RestAPIID == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.GetAuthorizerInput, error) {
	res := &svcsdk.GetAuthorizerInput{}

	if r.ko.Status.ID != nil {
		res.AuthorizerId = r.ko.Status.ID
	}
	if r.ko.Spec.RestAPIID != nil {
		res.RestApiId = r.ko.Spec.RestAPIID
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

	var resp *svcsdk.CreateAuthorizerOutput
	_ = resp
	resp, err = rm.sdkapi.CreateAuthorizer(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateAuthorizer", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.AuthType != nil {
		ko.Spec.AuthType = resp.AuthType
	} else {
		ko.Spec.AuthType = nil
	}
	if resp.AuthorizerCredentials != nil {
		ko.Spec.AuthorizerCredentials = resp.AuthorizerCredentials
	} else {
		ko.Spec.AuthorizerCredentials = nil
	}
	if resp.AuthorizerResultTtlInSeconds != nil {
		authorizerResultTTLInSecondsCopy := int64(*resp.AuthorizerResultTtlInSeconds)
		ko.Spec.AuthorizerResultTTLInSeconds = &authorizerResultTTLInSecondsCopy
	} else {
		ko.Spec.AuthorizerResultTTLInSeconds = nil
	}
	if resp.AuthorizerUri != nil {
		ko.Spec.AuthorizerURI = resp.AuthorizerUri
	} else {
		ko.Spec.AuthorizerURI = nil
	}
	if resp.Id != nil {
		ko.Status.ID = resp.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.IdentitySource != nil {
		ko.Spec.IdentitySource = resp.IdentitySource
	} else {
		ko.Spec.IdentitySource = nil
	}
	if resp.IdentityValidationExpression != nil {
		ko.Spec.IdentityValidationExpression = resp.IdentityValidationExpression
	} else {
		ko.Spec.IdentityValidationExpression = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ProviderARNs != nil {
		ko.Spec.ProviderARNs = aws.StringSlice(resp.ProviderARNs)
	} else {
		ko.Spec.ProviderARNs = nil
	}
	if resp.Type != "" {
		ko.Spec.Type = aws.String(string(resp.Type))
	} else {
		ko.Spec.Type = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateAuthorizerInput, error) {
	res := &svcsdk.CreateAuthorizerInput{}

	if r.ko.Spec.AuthType != nil {
		res.AuthType = r.ko.Spec.AuthType
	}
	if r.ko.Spec.AuthorizerCredentials != nil {
		res.AuthorizerCredentials = r.ko.Spec.AuthorizerCredentials
	}
	if r.ko.Spec.AuthorizerResultTTLInSeconds != nil {
		authorizerResultTTLInSecondsCopy0 := *r.ko.Spec.AuthorizerResultTTLInSeconds
		if authorizerResultTTLInSecondsCopy0 > math.MaxInt32 || authorizerResultTTLInSecondsCopy0 < math.MinInt32 {
			return nil, fmt.Errorf("error: field authorizerResultTtlInSeconds is of type int32")
		}
		authorizerResultTTLInSecondsCopy := int32(authorizerResultTTLInSecondsCopy0)
		res.AuthorizerResultTtlInSeconds = &authorizerResultTTLInSecondsCopy
	}
	if r.ko.Spec.AuthorizerURI != nil {
		res.AuthorizerUri = r.ko.Spec.AuthorizerURI
	}
	if r.ko.Spec.IdentitySource != nil {
		res.IdentitySource = r.ko.Spec.IdentitySource
	}
	if r.ko.Spec.IdentityValidationExpression != nil {
		res.IdentityValidationExpression = r.ko.Spec.IdentityValidationExpression
	}
	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}
	if r.ko.Spec.ProviderARNs != nil {
		res.ProviderARNs = aws.ToStringSlice(r.ko.Spec.ProviderARNs)
	}
	if r.ko.Spec.RestAPIID != nil {
		res.RestApiId = r.ko.Spec.RestAPIID
	}
	if r.ko.Spec.Type != nil {
		res.Type = svcsdktypes.AuthorizerType(*r.ko.Spec.Type)
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
	input, err := rm.newUpdateRequestPayload(ctx, desired, delta)
	if err != nil {
		return nil, err
	}
	updateAuthorizerInput(desired, latest, input, delta)

	var resp *svcsdk.UpdateAuthorizerOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateAuthorizer(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateAuthorizer", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.AuthType != nil {
		ko.Spec.AuthType = resp.AuthType
	} else {
		ko.Spec.AuthType = nil
	}
	if resp.AuthorizerCredentials != nil {
		ko.Spec.AuthorizerCredentials = resp.AuthorizerCredentials
	} else {
		ko.Spec.AuthorizerCredentials = nil
	}
	if resp.AuthorizerResultTtlInSeconds != nil {
		authorizerResultTTLInSecondsCopy := int64(*resp.AuthorizerResultTtlInSeconds)
		ko.Spec.AuthorizerResultTTLInSeconds = &authorizerResultTTLInSecondsCopy
	} else {
		ko.Spec.AuthorizerResultTTLInSeconds = nil
	}
	if resp.AuthorizerUri != nil {
		ko.Spec.AuthorizerURI = resp.AuthorizerUri
	} else {
		ko.Spec.AuthorizerURI = nil
	}
	if resp.Id != nil {
		ko.Status.ID = resp.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.IdentitySource != nil {
		ko.Spec.IdentitySource = resp.IdentitySource
	} else {
		ko.Spec.IdentitySource = nil
	}
	if resp.IdentityValidationExpression != nil {
		ko.Spec.IdentityValidationExpression = resp.IdentityValidationExpression
	} else {
		ko.Spec.IdentityValidationExpression = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ProviderARNs != nil {
		ko.Spec.ProviderARNs = aws.StringSlice(resp.ProviderARNs)
	} else {
		ko.Spec.ProviderARNs = nil
	}
	if resp.Type != "" {
		ko.Spec.Type = aws.String(string(resp.Type))
	} else {
		ko.Spec.Type = nil
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
) (*svcsdk.UpdateAuthorizerInput, error) {
	res := &svcsdk.UpdateAuthorizerInput{}

	if r.ko.Status.ID != nil {
		res.AuthorizerId = r.ko.Status.ID
	}
	if r.ko.Spec.RestAPIID != nil {
		res.RestApiId = r.ko.Spec.RestAPIID
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
	var resp *svcsdk.DeleteAuthorizerOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteAuthorizer(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteAuthorizer", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteAuthorizerInput, error) {
	res := &svcsdk.DeleteAuthorizerInput{}

	if r.ko.Status.ID != nil {
		res.AuthorizerId = r.ko.Status.ID
	}
	if r.ko.Spec.RestAPIID != nil {
		res.RestApiId = r.ko.Spec.RestAPIID
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Authorizer,
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
		"InvalidParameter":
		return true
	default:
		return false
	}
}
