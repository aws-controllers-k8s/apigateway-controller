	if delta.DifferentAt("Spec.Tags") {
		resourceARN := arnForResource(desired.ko, string(rm.awsPartition))
		if err := syncTags(ctx, rm.sdkapi, rm.metrics, resourceARN, desired.ko.Spec.Tags, latest.ko.Spec.Tags); err != nil {
			return nil, err
		}
	}
	if !delta.DifferentExcept("Spec.Tags") {
		return desired, nil
	}
