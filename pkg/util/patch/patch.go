package patch

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

var patchKeyEncoder = strings.NewReplacer("~", "~0", "/", "~1")

// A Set allows creating a set of patch operations that can be applied to a resource.
type Set struct {
	patchOps []*apigateway.PatchOperation
}

// ForSlice adds patch operations to this set at the specified path for adding new values and removing old values.
func (p *Set) ForSlice(path string, currValues, desiredValues []*string) {
	current := aws.StringValueSlice(currValues)
	desired := aws.StringValueSlice(desiredValues)
	var patchOps []*apigateway.PatchOperation
	for _, val := range current {
		if !slices.Contains(desired, val) {
			patchOps = append(patchOps, &apigateway.PatchOperation{
				Op:   aws.String(apigateway.OpRemove),
				Path: aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(val))),
			})
		}
	}
	for _, val := range desired {
		if !slices.Contains(current, val) {
			patchOps = append(patchOps, &apigateway.PatchOperation{
				Op:   aws.String(apigateway.OpAdd),
				Path: aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(val))),
			})
		}
	}
	p.patchOps = append(p.patchOps, patchOps...)
}

// ForMap adds patch operations to this set at the specified path for replacing existing keys with new values and
// removing keys that no longer exist.
func (p *Set) ForMap(path string, currValues, desiredValues map[string]*string, addSupported bool) {
	var patchOps []*apigateway.PatchOperation
	for k := range currValues {
		if _, ok := desiredValues[k]; !ok {
			patchOps = append(patchOps, &apigateway.PatchOperation{
				Op:   aws.String(apigateway.OpRemove),
				Path: aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(k))),
			})
		}
	}
	for k, v := range desiredValues {
		op := apigateway.OpReplace
		if addSupported {
			if _, ok := currValues[k]; !ok {
				op = apigateway.OpAdd
			}
		}
		patchOps = append(patchOps, &apigateway.PatchOperation{
			Op:    aws.String(op),
			Path:  aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(k))),
			Value: v,
		})
	}
	p.patchOps = append(p.patchOps, patchOps...)
}

// Replace adds a patch operation to this set for replacing the specified path with desiredVal.
func (p *Set) Replace(path string, desiredVal *string) {
	p.patchOps = append(p.patchOps, &apigateway.PatchOperation{
		Op:    aws.String(apigateway.OpReplace),
		Path:  aws.String(path),
		Value: desiredVal,
	})
}

// Remove adds a patch operation to this set for removing the specified path.
func (p *Set) Remove(path string) {
	p.patchOps = append(p.patchOps, &apigateway.PatchOperation{
		Op:   aws.String(apigateway.OpRemove),
		Path: aws.String(path),
	})
}

// GetPatchOperations returns the patch operations applied to this set.
func (p *Set) GetPatchOperations() []*apigateway.PatchOperation {
	if len(p.patchOps) > 0 {
		return p.patchOps
	}
	return nil
}
