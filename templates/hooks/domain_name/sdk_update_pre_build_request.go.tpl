	if delta.DifferentAt("Spec.Tags") {
		if err := rm.syncTags(ctx, latest, desired); err != nil {
			return nil, err
		}
	}
	if !delta.DifferentExcept("Spec.Tags") {
		return desired, nil
	}