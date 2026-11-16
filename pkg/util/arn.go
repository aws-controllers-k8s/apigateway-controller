package util

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

// apiGatewayServiceName is the service namespace used in API Gateway ARNs.
const apiGatewayServiceName = "apigateway"

// ARNForResource creates an ARN for the specified API Gateway resource. The
// partition is resolved by the resource manager from the aws-sdk-go-v2 client
// configuration and passed in by the caller.
func ARNForResource(resourceMeta *ackv1alpha1.ResourceMetadata, partition string, resourcePath string) string {
	return arn.ARN{
		Partition: partition,
		Service:   apiGatewayServiceName,
		Region:    string(*resourceMeta.Region),
		Resource:  resourcePath,
	}.String()
}
