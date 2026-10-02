	// GetIntegration returns the integration's HTTP method in httpMethod.
	// Only track it when the user declared one, so integrations that leave it
	// unset (for example MOCK) do not report drift.
	if ko.Spec.IntegrationHTTPMethod != nil && resp.HttpMethod != nil {
		ko.Spec.IntegrationHTTPMethod = resp.HttpMethod
	}
