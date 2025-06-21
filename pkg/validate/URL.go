package validate

import (
	"net/url"
)

func URL(URL string) bool {

	parsedURL, err := url.Parse(URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}
