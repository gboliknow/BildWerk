package utility

import "net/url"

func IsValidURL(link string) bool {
	parsedURL, err := url.ParseRequestURI(link)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}