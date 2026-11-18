	// GetIntegration returns the integration's HTTP method in httpMethod,
	// not the method request's. Spec.HTTPMethod identifies the method request
	// and is left as declared. MOCK integrations return no httpMethod, so
	// keep the declared value rather than reporting drift.
	if resp.HttpMethod != nil {
		ko.Spec.IntegrationHTTPMethod = resp.HttpMethod
	}
