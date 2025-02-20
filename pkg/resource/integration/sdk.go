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

package integration

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
	_ = &svcapitypes.Integration{}
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

	var resp *svcsdk.GetIntegrationOutput
	resp, err = rm.sdkapi.GetIntegration(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "GetIntegration", err)
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

	if resp.CacheKeyParameters != nil {
		ko.Spec.CacheKeyParameters = aws.StringSlice(resp.CacheKeyParameters)
	} else {
		ko.Spec.CacheKeyParameters = nil
	}
	if resp.CacheNamespace != nil {
		ko.Spec.CacheNamespace = resp.CacheNamespace
	} else {
		ko.Spec.CacheNamespace = nil
	}
	if resp.ConnectionId != nil {
		ko.Spec.ConnectionID = resp.ConnectionId
	} else {
		ko.Spec.ConnectionID = nil
	}
	if resp.ConnectionType != "" {
		ko.Spec.ConnectionType = aws.String(string(resp.ConnectionType))
	} else {
		ko.Spec.ConnectionType = nil
	}
	if resp.ContentHandling != "" {
		ko.Spec.ContentHandling = aws.String(string(resp.ContentHandling))
	} else {
		ko.Spec.ContentHandling = nil
	}
	if resp.Credentials != nil {
		ko.Spec.Credentials = resp.Credentials
	} else {
		ko.Spec.Credentials = nil
	}
	if resp.HttpMethod != nil {
		ko.Spec.HTTPMethod = resp.HttpMethod
	} else {
		ko.Spec.HTTPMethod = nil
	}
	if resp.IntegrationResponses != nil {
		f7 := map[string]*svcapitypes.IntegrationResponse{}
		for f7key, f7valiter := range resp.IntegrationResponses {
			f7val := &svcapitypes.IntegrationResponse{}
			if f7valiter.ContentHandling != "" {
				f7val.ContentHandling = aws.String(string(f7valiter.ContentHandling))
			}
			if f7valiter.ResponseParameters != nil {
				f7val.ResponseParameters = aws.StringMap(f7valiter.ResponseParameters)
			}
			if f7valiter.ResponseTemplates != nil {
				f7val.ResponseTemplates = aws.StringMap(f7valiter.ResponseTemplates)
			}
			if f7valiter.SelectionPattern != nil {
				f7val.SelectionPattern = f7valiter.SelectionPattern
			}
			if f7valiter.StatusCode != nil {
				f7val.StatusCode = f7valiter.StatusCode
			}
			f7[f7key] = f7val
		}
		ko.Status.IntegrationResponses = f7
	} else {
		ko.Status.IntegrationResponses = nil
	}
	if resp.PassthroughBehavior != nil {
		ko.Spec.PassthroughBehavior = resp.PassthroughBehavior
	} else {
		ko.Spec.PassthroughBehavior = nil
	}
	if resp.RequestParameters != nil {
		ko.Spec.RequestParameters = aws.StringMap(resp.RequestParameters)
	} else {
		ko.Spec.RequestParameters = nil
	}
	if resp.RequestTemplates != nil {
		ko.Spec.RequestTemplates = aws.StringMap(resp.RequestTemplates)
	} else {
		ko.Spec.RequestTemplates = nil
	}
	timeoutInMillisCopy := int64(resp.TimeoutInMillis)
	ko.Spec.TimeoutInMillis = &timeoutInMillisCopy
	if resp.TlsConfig != nil {
		f12 := &svcapitypes.TLSConfig{}
		f12.InsecureSkipVerification = &resp.TlsConfig.InsecureSkipVerification
		ko.Spec.TLSConfig = f12
	} else {
		ko.Spec.TLSConfig = nil
	}
	if resp.Type != "" {
		ko.Spec.Type = aws.String(string(resp.Type))
	} else {
		ko.Spec.Type = nil
	}
	if resp.Uri != nil {
		ko.Spec.URI = resp.Uri
	} else {
		ko.Spec.URI = nil
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
	return r.ko.Spec.HTTPMethod == nil || r.ko.Spec.ResourceID == nil || r.ko.Spec.RestAPIID == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.GetIntegrationInput, error) {
	res := &svcsdk.GetIntegrationInput{}

	if r.ko.Spec.HTTPMethod != nil {
		res.HttpMethod = r.ko.Spec.HTTPMethod
	}
	if r.ko.Spec.ResourceID != nil {
		res.ResourceId = r.ko.Spec.ResourceID
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

	var resp *svcsdk.PutIntegrationOutput
	_ = resp
	resp, err = rm.sdkapi.PutIntegration(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "PutIntegration", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.CacheKeyParameters != nil {
		ko.Spec.CacheKeyParameters = aws.StringSlice(resp.CacheKeyParameters)
	} else {
		ko.Spec.CacheKeyParameters = nil
	}
	if resp.CacheNamespace != nil {
		ko.Spec.CacheNamespace = resp.CacheNamespace
	} else {
		ko.Spec.CacheNamespace = nil
	}
	if resp.ConnectionId != nil {
		ko.Spec.ConnectionID = resp.ConnectionId
	} else {
		ko.Spec.ConnectionID = nil
	}
	if resp.ConnectionType != "" {
		ko.Spec.ConnectionType = aws.String(string(resp.ConnectionType))
	} else {
		ko.Spec.ConnectionType = nil
	}
	if resp.ContentHandling != "" {
		ko.Spec.ContentHandling = aws.String(string(resp.ContentHandling))
	} else {
		ko.Spec.ContentHandling = nil
	}
	if resp.Credentials != nil {
		ko.Spec.Credentials = resp.Credentials
	} else {
		ko.Spec.Credentials = nil
	}
	if resp.HttpMethod != nil {
		ko.Spec.HTTPMethod = resp.HttpMethod
	} else {
		ko.Spec.HTTPMethod = nil
	}
	if resp.IntegrationResponses != nil {
		f7 := map[string]*svcapitypes.IntegrationResponse{}
		for f7key, f7valiter := range resp.IntegrationResponses {
			f7val := &svcapitypes.IntegrationResponse{}
			if f7valiter.ContentHandling != "" {
				f7val.ContentHandling = aws.String(string(f7valiter.ContentHandling))
			}
			if f7valiter.ResponseParameters != nil {
				f7val.ResponseParameters = aws.StringMap(f7valiter.ResponseParameters)
			}
			if f7valiter.ResponseTemplates != nil {
				f7val.ResponseTemplates = aws.StringMap(f7valiter.ResponseTemplates)
			}
			if f7valiter.SelectionPattern != nil {
				f7val.SelectionPattern = f7valiter.SelectionPattern
			}
			if f7valiter.StatusCode != nil {
				f7val.StatusCode = f7valiter.StatusCode
			}
			f7[f7key] = f7val
		}
		ko.Status.IntegrationResponses = f7
	} else {
		ko.Status.IntegrationResponses = nil
	}
	if resp.PassthroughBehavior != nil {
		ko.Spec.PassthroughBehavior = resp.PassthroughBehavior
	} else {
		ko.Spec.PassthroughBehavior = nil
	}
	if resp.RequestParameters != nil {
		ko.Spec.RequestParameters = aws.StringMap(resp.RequestParameters)
	} else {
		ko.Spec.RequestParameters = nil
	}
	if resp.RequestTemplates != nil {
		ko.Spec.RequestTemplates = aws.StringMap(resp.RequestTemplates)
	} else {
		ko.Spec.RequestTemplates = nil
	}
	timeoutInMillisCopy := int64(resp.TimeoutInMillis)
	ko.Spec.TimeoutInMillis = &timeoutInMillisCopy
	if resp.TlsConfig != nil {
		f12 := &svcapitypes.TLSConfig{}
		f12.InsecureSkipVerification = &resp.TlsConfig.InsecureSkipVerification
		ko.Spec.TLSConfig = f12
	} else {
		ko.Spec.TLSConfig = nil
	}
	if resp.Type != "" {
		ko.Spec.Type = aws.String(string(resp.Type))
	} else {
		ko.Spec.Type = nil
	}
	if resp.Uri != nil {
		ko.Spec.URI = resp.Uri
	} else {
		ko.Spec.URI = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.PutIntegrationInput, error) {
	res := &svcsdk.PutIntegrationInput{}

	if r.ko.Spec.CacheKeyParameters != nil {
		res.CacheKeyParameters = aws.ToStringSlice(r.ko.Spec.CacheKeyParameters)
	}
	if r.ko.Spec.CacheNamespace != nil {
		res.CacheNamespace = r.ko.Spec.CacheNamespace
	}
	if r.ko.Spec.ConnectionID != nil {
		res.ConnectionId = r.ko.Spec.ConnectionID
	}
	if r.ko.Spec.ConnectionType != nil {
		res.ConnectionType = svcsdktypes.ConnectionType(*r.ko.Spec.ConnectionType)
	}
	if r.ko.Spec.ContentHandling != nil {
		res.ContentHandling = svcsdktypes.ContentHandlingStrategy(*r.ko.Spec.ContentHandling)
	}
	if r.ko.Spec.Credentials != nil {
		res.Credentials = r.ko.Spec.Credentials
	}
	if r.ko.Spec.HTTPMethod != nil {
		res.HttpMethod = r.ko.Spec.HTTPMethod
	}
	if r.ko.Spec.IntegrationHTTPMethod != nil {
		res.IntegrationHttpMethod = r.ko.Spec.IntegrationHTTPMethod
	}
	if r.ko.Spec.PassthroughBehavior != nil {
		res.PassthroughBehavior = r.ko.Spec.PassthroughBehavior
	}
	if r.ko.Spec.RequestParameters != nil {
		res.RequestParameters = aws.ToStringMap(r.ko.Spec.RequestParameters)
	}
	if r.ko.Spec.RequestTemplates != nil {
		res.RequestTemplates = aws.ToStringMap(r.ko.Spec.RequestTemplates)
	}
	if r.ko.Spec.ResourceID != nil {
		res.ResourceId = r.ko.Spec.ResourceID
	}
	if r.ko.Spec.RestAPIID != nil {
		res.RestApiId = r.ko.Spec.RestAPIID
	}
	if r.ko.Spec.TimeoutInMillis != nil {
		timeoutInMillisCopy0 := *r.ko.Spec.TimeoutInMillis
		if timeoutInMillisCopy0 > math.MaxInt32 || timeoutInMillisCopy0 < math.MinInt32 {
			return nil, fmt.Errorf("error: field TimeoutInMillis is of type int32")
		}
		timeoutInMillisCopy := int32(timeoutInMillisCopy0)
		res.TimeoutInMillis = &timeoutInMillisCopy
	}
	if r.ko.Spec.TLSConfig != nil {
		f14 := &svcsdktypes.TlsConfig{}
		if r.ko.Spec.TLSConfig.InsecureSkipVerification != nil {
			f14.InsecureSkipVerification = *r.ko.Spec.TLSConfig.InsecureSkipVerification
		}
		res.TlsConfig = f14
	}
	if r.ko.Spec.Type != nil {
		res.Type = svcsdktypes.IntegrationType(*r.ko.Spec.Type)
	}
	if r.ko.Spec.URI != nil {
		res.Uri = r.ko.Spec.URI
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
	updateIntegrationInput(desired, latest, input, delta)

	var resp *svcsdk.UpdateIntegrationOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateIntegration(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateIntegration", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.CacheKeyParameters != nil {
		ko.Spec.CacheKeyParameters = aws.StringSlice(resp.CacheKeyParameters)
	} else {
		ko.Spec.CacheKeyParameters = nil
	}
	if resp.CacheNamespace != nil {
		ko.Spec.CacheNamespace = resp.CacheNamespace
	} else {
		ko.Spec.CacheNamespace = nil
	}
	if resp.ConnectionId != nil {
		ko.Spec.ConnectionID = resp.ConnectionId
	} else {
		ko.Spec.ConnectionID = nil
	}
	if resp.ConnectionType != "" {
		ko.Spec.ConnectionType = aws.String(string(resp.ConnectionType))
	} else {
		ko.Spec.ConnectionType = nil
	}
	if resp.ContentHandling != "" {
		ko.Spec.ContentHandling = aws.String(string(resp.ContentHandling))
	} else {
		ko.Spec.ContentHandling = nil
	}
	if resp.Credentials != nil {
		ko.Spec.Credentials = resp.Credentials
	} else {
		ko.Spec.Credentials = nil
	}
	if resp.HttpMethod != nil {
		ko.Spec.HTTPMethod = resp.HttpMethod
	} else {
		ko.Spec.HTTPMethod = nil
	}
	if resp.IntegrationResponses != nil {
		f7 := map[string]*svcapitypes.IntegrationResponse{}
		for f7key, f7valiter := range resp.IntegrationResponses {
			f7val := &svcapitypes.IntegrationResponse{}
			if f7valiter.ContentHandling != "" {
				f7val.ContentHandling = aws.String(string(f7valiter.ContentHandling))
			}
			if f7valiter.ResponseParameters != nil {
				f7val.ResponseParameters = aws.StringMap(f7valiter.ResponseParameters)
			}
			if f7valiter.ResponseTemplates != nil {
				f7val.ResponseTemplates = aws.StringMap(f7valiter.ResponseTemplates)
			}
			if f7valiter.SelectionPattern != nil {
				f7val.SelectionPattern = f7valiter.SelectionPattern
			}
			if f7valiter.StatusCode != nil {
				f7val.StatusCode = f7valiter.StatusCode
			}
			f7[f7key] = f7val
		}
		ko.Status.IntegrationResponses = f7
	} else {
		ko.Status.IntegrationResponses = nil
	}
	if resp.PassthroughBehavior != nil {
		ko.Spec.PassthroughBehavior = resp.PassthroughBehavior
	} else {
		ko.Spec.PassthroughBehavior = nil
	}
	if resp.RequestParameters != nil {
		ko.Spec.RequestParameters = aws.StringMap(resp.RequestParameters)
	} else {
		ko.Spec.RequestParameters = nil
	}
	if resp.RequestTemplates != nil {
		ko.Spec.RequestTemplates = aws.StringMap(resp.RequestTemplates)
	} else {
		ko.Spec.RequestTemplates = nil
	}
	timeoutInMillisCopy := int64(resp.TimeoutInMillis)
	ko.Spec.TimeoutInMillis = &timeoutInMillisCopy
	if resp.TlsConfig != nil {
		f12 := &svcapitypes.TLSConfig{}
		f12.InsecureSkipVerification = &resp.TlsConfig.InsecureSkipVerification
		ko.Spec.TLSConfig = f12
	} else {
		ko.Spec.TLSConfig = nil
	}
	if resp.Type != "" {
		ko.Spec.Type = aws.String(string(resp.Type))
	} else {
		ko.Spec.Type = nil
	}
	if resp.Uri != nil {
		ko.Spec.URI = resp.Uri
	} else {
		ko.Spec.URI = nil
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
) (*svcsdk.UpdateIntegrationInput, error) {
	res := &svcsdk.UpdateIntegrationInput{}

	if r.ko.Spec.HTTPMethod != nil {
		res.HttpMethod = r.ko.Spec.HTTPMethod
	}
	if r.ko.Spec.ResourceID != nil {
		res.ResourceId = r.ko.Spec.ResourceID
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
	var resp *svcsdk.DeleteIntegrationOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteIntegration(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteIntegration", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteIntegrationInput, error) {
	res := &svcsdk.DeleteIntegrationInput{}

	if r.ko.Spec.HTTPMethod != nil {
		res.HttpMethod = r.ko.Spec.HTTPMethod
	}
	if r.ko.Spec.ResourceID != nil {
		res.ResourceId = r.ko.Spec.ResourceID
	}
	if r.ko.Spec.RestAPIID != nil {
		res.RestApiId = r.ko.Spec.RestAPIID
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Integration,
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
