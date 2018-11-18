package page

import (
	"strings"
	"sync"
	"testing"
)

func TestValidatePage(t *testing.T) {
	domain = "http://toto"
	testCases := []struct {
		id     int
		url    string
		assert string
	}{
		{1, "/", ""},
		{2, "test", "/test"},
		{3, "test/", "/test"},
		{4, "http://test/", ""},
		{5, "http://toto/", ""},
		{6, "http://toto/test", "/test"},
		{7, "http://toto/test/test2", "/test/test2"},
		{7, "http://toto/test/#", "/test"},
		{7, "http://toto/test/#toto", "/test"},
		{7, "http://toto/test/test2#toto", "/test/test2"},
	}

	for _, p := range testCases {
		b := simplifyURL(p.url)
		if b != p.assert {
			t.Errorf("Test %d for url %s should be of validity %s, assert %s\n", p.id, p.url, p.assert, b)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	pages.links = map[string]struct{}{
		"/test/":       struct{}{},
		"/test/test2/": struct{}{},
	}
	testCases := []struct {
		id     int
		links  []string
		assert []string
	}{
		{1, []string{"/test/"}, []string{}},
		{2, []string{"/test/", "/toto/"}, []string{"/toto/"}},
		{3, []string{"/toto/"}, []string{}}, // Added to the global var pages
		{4, []string{"/toto/", "/tata/"}, []string{"/tata/"}},
		{5, []string{"/toto/", "/tata/"}, []string{}},
		{6, []string{"/titi/", "/tutu/"}, []string{"/titi/", "/tutu/"}},
	}

	var mux = sync.Mutex{}
	for _, tc := range testCases {
		d := removeDuplicateLinks(tc.links, &mux)
		tcStr := strings.Join(tc.assert[:], ",")
		dStr := strings.Join(d[:], ",")

		if tcStr != dStr {
			t.Errorf("Test %d - %v should be of validity %s, assert %s\n", tc.id, tc.links, tc.assert, d)
		}
	}

	if len(pages.links) == 2 {
		t.Error("Elements were not added to global variable \"pages\"")
	}
}
