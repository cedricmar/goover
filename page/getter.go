package page

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/cedricmar/goover/helper"
	"golang.org/x/net/html"
)

type pageResources struct {
	links, files map[string]struct{}
}

var (
	domain string
	pages  = pageResources{
		links: map[string]struct{}{},
		files: map[string]struct{}{},
	}
)

// Get webpage
func Get(dom string, url string, dest string, wg *sync.WaitGroup, mux *sync.Mutex) {
	// Set the domain
	if domain == "" {
		domain = dom
	}

	// Fetch the page
	resp, err := http.Get(domain + url)
	helper.Check(err)
	defer resp.Body.Close()

	// The page was found
	if resp.StatusCode == http.StatusOK {
		// Save page
		b, err := ioutil.ReadAll(resp.Body)
		helper.Check(err)
		wg.Add(1)
		go savePage(dest, url, b, wg)

		// Reset response read state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		res := getPageResources(resp)

		// Download files
		for _, l := range removeDuplicateFiles(res["files"], mux) {
			fmt.Println("Fetching " + domain + l)
			// Save
			wg.Add(1)
			go saveResource(dest, l, wg)
		}

		// Follow links
		for _, l := range removeDuplicateLinks(res["links"], mux) {
			fmt.Println("Fetching " + domain + l)
			wg.Add(1)
			go Get(domain, l, dest, wg, mux)
		}
	}

	wg.Done()
}

// Analyze page and return
func getPageResources(r *http.Response) map[string][]string {

	// List of resources to save
	res := make(map[string][]string)
	res["links"] = []string{}
	res["files"] = []string{}

	if r.StatusCode != http.StatusOK {
		return res
	}

	z := html.NewTokenizer(r.Body)

	for {
		tt := z.Next()

		switch {
		// End of the document
		case tt == html.ErrorToken:
			return res
		// Tags
		case tt == html.StartTagToken:
			t := z.Token()

			switch {
			// Links
			case t.Data == "a":
				url := getHREF(&t)
				if url = simplifyURL(url); url != "" {
					res["links"] = append(res["links"], url)
				}
			case t.Data == "script" || t.Data == "img":
				url := getSRC(&t)
				if url = simplifyURL(url); url != "" {
					res["files"] = append(res["files"], url)
				}
			case t.Data == "style" || t.Data == "link":
				url := getHREF(&t)
				if url = simplifyURL(url); url != "" {
					res["files"] = append(res["files"], url)
				}
			}
		}
	}
}

func getHREF(t *html.Token) string {
	return getATTR("href", t)
}

func getSRC(t *html.Token) string {
	return getATTR("src", t)
}

func getATTR(name string, t *html.Token) string {
	for _, a := range t.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}
