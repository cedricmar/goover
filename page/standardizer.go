package page

import (
	"strings"
	"sync"
)

func simplifyURL(url string) string {
	// External link
	if (strings.Contains(url, "http://") || strings.Contains(url, "https://")) && !strings.Contains(url, domain) {
		return ""
	}
	// Remove domain for the rest
	url = strings.Replace(url, domain, "", -1)
	// Remove #anchors
	if strings.Contains(url, "#") {
		parts := strings.Split(url, "#")
		url = parts[0]
	}
	// Home
	if url == "/" || url == "" {
		return ""
	}
	// Standardizing
	return "/" + strings.Trim(url, "/")
}

func removeDuplicateLinks(urls []string, mux *sync.Mutex) []string {
	noDupes := []string{}
	for _, l := range urls {
		mux.Lock()
		if _, ok := pages.links[l]; !ok {
			// Does not exist
			pages.links[l] = struct{}{}
			noDupes = append(noDupes, l)
		}
		mux.Unlock()
	}
	return noDupes
}

func removeDuplicateFiles(urls []string, mux *sync.Mutex) []string {
	noDupes := []string{}
	for _, l := range urls {
		mux.Lock()
		if _, ok := pages.files[l]; !ok {
			// Does not exist
			pages.files[l] = struct{}{}
			noDupes = append(noDupes, l)
		}
		mux.Unlock()
	}
	return noDupes
}
