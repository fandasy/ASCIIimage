package validate

func ContentType(contentType string, types ...string) bool {
	for _, v := range types {
		if contentType == v {
			return true
		}
	}

	return false
}
