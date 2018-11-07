package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

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

func usage() {
	fmt.Printf("Usage:\n  goover <url> <folder>\n")
	os.Exit(0)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
	//path := "./" + dest
	//filename := "index.html"
	//body, _ := ioutil.ReadAll(resp.Body)

	// Save this page
	//save(path, filename, body)
*/
func save(path string, filename string, txt []byte) {
	_ = os.Mkdir(path, os.ModePerm)
	err := ioutil.WriteFile(path+filename, txt, os.ModePerm)
	check(err)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	domain = os.Args[1]
	dest := os.Args[2]

	// @TODO - sanity check for url and folder
	// @TODO - create folder if not exist

	fmt.Printf("Hoovering %s in %s\n", domain, dest)

	var wg sync.WaitGroup
	var mux = sync.Mutex{}

	wg.Add(1)
	go getPage("/", dest, &wg, &mux)

	wg.Wait()

	fmt.Println(len(pages.links), pages)

	fmt.Println("Done")
}

// @TODO - implement worker pool for big websites maybe ?
func getPage(url string, dest string, wg *sync.WaitGroup, mux *sync.Mutex) {

	// Fetch the page
	resp, err := http.Get(domain + url)
	check(err)
	defer resp.Body.Close()

	// The page was found
	if resp.StatusCode == http.StatusOK {

		// Save page
		//wg.Add(1)

		res := getPageResources(resp)

		fmt.Println(res)
		// Download files
		for _, l := range removeDuplicateFiles(res["files"], mux) {
			fmt.Println("Fetching " + domain + l)
			// Save
			//wg.Add(1)
		}

		// Follow links
		for _, l := range removeDuplicateLinks(res["links"], mux) {
			fmt.Println("Fetching " + domain + l)
			wg.Add(1)
			go getPage(l, dest, wg, mux)
		}
	}

	wg.Done()
}

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
