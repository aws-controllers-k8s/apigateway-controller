package patch

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	apigatewaytypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

var patchKeyEncoder = strings.NewReplacer("~", "~0", "/", "~1")

// A Set allows creating a set of patch operations that can be applied to a resource.
type Set struct {
	patchOps []*apigatewaytypes.PatchOperation
}

// ForSlice adds patch operations to this set at the specified path for adding new values and removing old values.
func (p *Set) ForSlice(path string, currValues, desiredValues []*string) {
	current := aws.ToStringSlice(currValues)
	desired := aws.ToStringSlice(desiredValues)
	var patchOps []*apigatewaytypes.PatchOperation
	for _, val := range current {
		if !slices.Contains(desired, val) {
			patchOps = append(patchOps, &apigatewaytypes.PatchOperation{
				Op:   apigatewaytypes.OpRemove,
				Path: aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(val))),
			})
		}
	}
	for _, val := range desired {
		if !slices.Contains(current, val) {
			patchOps = append(patchOps, &apigatewaytypes.PatchOperation{
				Op:   apigatewaytypes.OpAdd,
				Path: aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(val))),
			})
		}
	}
	p.patchOps = append(p.patchOps, patchOps...)
}

// ForMap adds patch operations to this set at the specified path for replacing existing keys with new values and
// removing keys that no longer exist.
func (p *Set) ForMap(path string, currValues, desiredValues map[string]*string, addSupported bool) {
	var patchOps []*apigatewaytypes.PatchOperation
	for k := range currValues {
		if _, ok := desiredValues[k]; !ok {
			patchOps = append(patchOps, &apigatewaytypes.PatchOperation{
				Op:   apigatewaytypes.OpRemove,
				Path: aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(k))),
			})
		}
	}
	for k, v := range desiredValues {
		op := apigatewaytypes.OpReplace
		if addSupported {
			if _, ok := currValues[k]; !ok {
				op = apigatewaytypes.OpAdd
			}
		}
		patchOps = append(patchOps, &apigatewaytypes.PatchOperation{
			Op:    op,
			Path:  aws.String(fmt.Sprintf("%s/%s", path, patchKeyEncoder.Replace(k))),
			Value: v,
		})
	}
	p.patchOps = append(p.patchOps, patchOps...)
}

// Replace adds a patch operation to this set for replacing the specified path with desiredVal.
func (p *Set) Replace(path string, desiredVal *string) {
	p.patchOps = append(p.patchOps, &apigatewaytypes.PatchOperation{
		Op:    apigatewaytypes.OpReplace,
		Path:  aws.String(path),
		Value: desiredVal,
	})
}

// Remove adds a patch operation to this set for removing the specified path.
func (p *Set) Remove(path string) {
	p.patchOps = append(p.patchOps, &apigatewaytypes.PatchOperation{
		Op:   apigatewaytypes.OpRemove,
		Path: aws.String(path),
	})
}

// GetPatchOperations returns the patch operations applied to this set.
func (p *Set) GetPatchOperations() []apigatewaytypes.PatchOperation {
	var patchOps []apigatewaytypes.PatchOperation
	for _, op := range p.patchOps {
		patchOps = append(patchOps, *op)
	}
	return patchOps
}
