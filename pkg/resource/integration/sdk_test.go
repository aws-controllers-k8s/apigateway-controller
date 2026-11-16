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

package integration

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/apigateway-controller/apis/v1alpha1"
)

// fakeHTTPClient returns a canned GetIntegration response body and records
// the request path so tests can check which method request was read.
type fakeHTTPClient struct {
	body        string
	requestPath string
}

func (c *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.requestPath = req.URL.Path
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(c.body)),
		Request:    req,
	}, nil
}

func newTestResourceManager(httpClient *fakeHTTPClient) *resourceManager {
	cfg := aws.Config{
		Region:      "us-west-2",
		Credentials: aws.AnonymousCredentials{},
		HTTPClient:  httpClient,
	}
	return &resourceManager{
		metrics: ackmetrics.NewMetrics("apigateway"),
		sdkapi:  svcsdk.NewFromConfig(cfg),
	}
}

func TestSdkFindIntegrationHTTPMethod(t *testing.T) {
	for _, tt := range []struct {
		description                   string
		desired                       svcapitypes.IntegrationSpec
		responseBody                  string
		expectedHTTPMethod            *string
		expectedIntegrationHTTPMethod *string
	}{
		{
			description: "AWS_PROXY with distinct request and backend methods",
			desired: svcapitypes.IntegrationSpec{
				HTTPMethod:            aws.String("GET"),
				IntegrationHTTPMethod: aws.String("POST"),
				Type:                  aws.String("AWS_PROXY"),
			},
			responseBody:                  `{"type":"AWS_PROXY","httpMethod":"POST"}`,
			expectedHTTPMethod:            aws.String("GET"),
			expectedIntegrationHTTPMethod: aws.String("POST"),
		},
		{
			description: "backend method drift is observed",
			desired: svcapitypes.IntegrationSpec{
				HTTPMethod:            aws.String("GET"),
				IntegrationHTTPMethod: aws.String("POST"),
				Type:                  aws.String("AWS_PROXY"),
			},
			responseBody:                  `{"type":"AWS_PROXY","httpMethod":"GET"}`,
			expectedHTTPMethod:            aws.String("GET"),
			expectedIntegrationHTTPMethod: aws.String("GET"),
		},
		{
			description: "adoption without integrationHTTPMethod reads backend method",
			desired: svcapitypes.IntegrationSpec{
				HTTPMethod: aws.String("GET"),
			},
			responseBody:                  `{"type":"AWS_PROXY","httpMethod":"POST"}`,
			expectedHTTPMethod:            aws.String("GET"),
			expectedIntegrationHTTPMethod: aws.String("POST"),
		},
		{
			description: "MOCK with omitted backend method keeps request method",
			desired: svcapitypes.IntegrationSpec{
				HTTPMethod: aws.String("OPTIONS"),
				Type:       aws.String("MOCK"),
			},
			responseBody:                  `{"type":"MOCK"}`,
			expectedHTTPMethod:            aws.String("OPTIONS"),
			expectedIntegrationHTTPMethod: nil,
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			httpClient := &fakeHTTPClient{body: tt.responseBody}
			rm := newTestResourceManager(httpClient)

			spec := tt.desired
			spec.RestAPIID = aws.String("api-id")
			spec.ResourceID = aws.String("resource-id")
			desired := &resource{ko: &svcapitypes.Integration{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
				Spec:       spec,
			}}

			latest, err := rm.sdkFind(context.Background(), desired)
			require.NoError(t, err)

			assert.Equal(t,
				"/restapis/api-id/resources/resource-id/methods/"+*tt.desired.HTTPMethod+"/integration",
				httpClient.requestPath)
			assert.Equal(t, tt.expectedHTTPMethod, latest.ko.Spec.HTTPMethod)
			assert.Equal(t, tt.expectedIntegrationHTTPMethod, latest.ko.Spec.IntegrationHTTPMethod)
		})
	}
}
