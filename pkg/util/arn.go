package util

import (
	"fmt"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

// ARNForResource creates an ARN for the specified API Gateway resource.
func ARNForResource(resourceMeta *ackv1alpha1.ResourceMetadata, resourcePath string) (string, error) {
	region := string(*resourceMeta.Region)
	partition, found := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region)
	if !found {
		return "", fmt.Errorf("failed to find partition for region %q", region)
	}

	return arn.ARN{
		Partition: partition.ID(),
		Service:   apigateway.ServiceName,
		Region:    region,
		Resource:  resourcePath,
	}.String(), nil
}
