package page

import (
	"io/ioutil"
	"os"

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
func Save(path string, filename string, txt []byte) {
	_ = os.Mkdir(path, os.ModePerm)
	err := ioutil.WriteFile(path+filename, txt, os.ModePerm)
	helper.Check(err)
}
