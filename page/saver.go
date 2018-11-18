package page

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/cedricmar/goover/helper"
)

/*
	//path := "./" + dest
	//filename := "index.html"
	//body, _ := ioutil.ReadAll(resp.Body)

	// Save this page
	//save(path, filename, body)
*/

// Save files on disk
func savePage(dest string, path string, txt []byte, wg *sync.WaitGroup) {

	// Standardize filepath
	savePath := getSavePath(dest, path)

	_ = os.Mkdir(savePath, os.ModePerm)
	err := ioutil.WriteFile(savePath+"index.html", txt, os.ModePerm)
	helper.Check(err)

	wg.Done()
}

func saveResource(dest string, url string, wg *sync.WaitGroup) {
	// Download it
	resp, err := http.Get(domain + url)
	helper.Check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		res, err := ioutil.ReadAll(resp.Body)
		helper.Check(err)

		filename, savePath := getFilePath(dest, simplifyURL(url))

		fmt.Println(">>> savepath - " + savePath)
		fmt.Println(">>> filename - " + filename)

		_ = os.Mkdir(savePath, os.ModePerm)
		err = ioutil.WriteFile(savePath+filename, res, os.ModePerm)
		helper.Check(err)
	}

	wg.Done()
}
