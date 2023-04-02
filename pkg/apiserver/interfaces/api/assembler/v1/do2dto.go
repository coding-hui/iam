package v1

func convertBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
