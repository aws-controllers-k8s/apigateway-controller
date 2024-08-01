package patch_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/stretchr/testify/assert"

	"github.com/aws-controllers-k8s/apigateway-controller/pkg/util/patch"
)

func TestPatchOperations(t *testing.T) {
	for _, tt := range []struct {
		applyPatches func(*patch.Set)
		description  string

		expectedPatchOps []*apigateway.PatchOperation
	}{
		{
			description: "all supported patch operations",
			applyPatches: func(patchSet *patch.Set) {
				patchSet.Replace("/field", aws.String("newValue"))
				patchSet.ForSlice("/items", aws.StringSlice([]string{"a", "b", "c"}), aws.StringSlice([]string{"b", "d"}))
				patchSet.ForMap("/keys", map[string]*string{
					"k1": aws.String("v1"),
					"k2": aws.String("v2"),
				}, map[string]*string{
					"k2": aws.String("v5"),
					"k1": aws.String("v1"),
					"k3": aws.String("v3"),
				}, false)
				patchSet.Remove("/removed")
			},
			expectedPatchOps: []*apigateway.PatchOperation{
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/field"),
					Value: aws.String("newValue"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/items/a"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/items/c"),
				},
				{
					Op:   aws.String(apigateway.OpAdd),
					Path: aws.String("/items/d"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k2"),
					Value: aws.String("v5"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k1"),
					Value: aws.String("v1"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k3"),
					Value: aws.String("v3"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/removed"),
				},
			},
		},
		{
			description: "all supported patch operations with keywords",
			applyPatches: func(patchSet *patch.Set) {
				patchSet.Replace("/field", aws.String("~newValue"))
				patchSet.ForSlice("/items", aws.StringSlice([]string{"/a", "~b", "c"}), aws.StringSlice([]string{"~b", "~d"}))
				patchSet.ForMap("/keys", map[string]*string{
					"k1~":  aws.String("v1~/"),
					"k2/~": aws.String("v2~/"),
				}, map[string]*string{
					"k2/~": aws.String("v5~/"),
					"k1~":  aws.String("v1~/"),
					"k3~/": aws.String("v3~/"),
				}, false)
				patchSet.Remove("/removed")
			},
			expectedPatchOps: []*apigateway.PatchOperation{
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/field"),
					Value: aws.String("~newValue"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/items/~1a"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/items/c"),
				},
				{
					Op:   aws.String(apigateway.OpAdd),
					Path: aws.String("/items/~0d"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k2~1~0"),
					Value: aws.String("v5~/"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k1~0"),
					Value: aws.String("v1~/"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k3~0~1"),
					Value: aws.String("v3~/"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/removed"),
				},
			},
		},
		{
			description: "maps support OpAdd",
			applyPatches: func(patchSet *patch.Set) {
				patchSet.Replace("/field", aws.String("newValue"))
				patchSet.ForSlice("/items", aws.StringSlice([]string{"a", "b", "c"}), aws.StringSlice([]string{"b", "d"}))
				patchSet.ForMap("/keys", map[string]*string{
					"k1": aws.String("v1"),
					"k2": aws.String("v2"),
				}, map[string]*string{
					"k2": aws.String("v5"),
					"k1": aws.String("v1"),
					"k3": aws.String("v3"),
				}, true)
				patchSet.Remove("/removed")
			},
			expectedPatchOps: []*apigateway.PatchOperation{
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/field"),
					Value: aws.String("newValue"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/items/a"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/items/c"),
				},
				{
					Op:   aws.String(apigateway.OpAdd),
					Path: aws.String("/items/d"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k2"),
					Value: aws.String("v5"),
				},
				{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String("/keys/k1"),
					Value: aws.String("v1"),
				},
				{
					Op:    aws.String(apigateway.OpAdd),
					Path:  aws.String("/keys/k3"),
					Value: aws.String("v3"),
				},
				{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String("/removed"),
				},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			assert := assert.New(t)
			var patchSet patch.Set
			tt.applyPatches(&patchSet)
			assert.ElementsMatch(patchSet.GetPatchOperations(), tt.expectedPatchOps)
		})
	}
}
