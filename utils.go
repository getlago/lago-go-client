package lago

func Ptr[T any](v T) *T {
	return &v
}

// statusQueryParams builds a query param map from an optional status variadic arg.
func statusQueryParams(subscriptionStatus []string) map[string]string {
	if len(subscriptionStatus) > 0 && subscriptionStatus[0] != "" {
		return map[string]string{"subscription_status": subscriptionStatus[0]}
	}
	return nil
}
