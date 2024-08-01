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
