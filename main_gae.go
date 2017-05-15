// +build appengine

package gitignore

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	fileList  []string
	ignoreMap map[string][]byte
	homePage  []byte
)

func init() {
	r := httprouter.New()

	r.ServeFiles("/public/*filepath", http.Dir("public/static/"))
	r.GET("/api/:name", apiHandler)
	r.GET("/", mainHandler)

	http.Handle("/", r)
}

func init() {
	ignoreMap = initIgnoreFiles()
	homePage = initIndexPage()
}

func mainHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write(homePage)
}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	params := ps.ByName("name")
	names := strings.Split(params, ",")
	sort.Strings(names)
	res := apiResponse(names)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(res)
}

func apiResponse(names []string) (res []byte) {
	for _, name := range names {
		ignoreTitle := []byte(fmt.Sprintf("\n%s %s %s\n", "###", strings.Title(name), "###"))
		content := ignoreMap[name]
		content = append(ignoreTitle, content...)
		res = append(res, content...)
	}

	return
}

func initIndexPage() (res []byte) {
	res, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		log.Fatalf("Error, parsing index.html: %v", err)
	}

	return
}

func initIgnoreFiles() (list map[string][]byte) {
	list = make(map[string][]byte)

	err := filepath.Walk("data/gitignore/", walker)
	if err != nil {
		log.Fatalf("Error, in filepath.Walk(): %v", err)
	}

	for _, file := range fileList {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Error, reading files: %v", err)
		}
		name := strings.ToLower(strings.TrimSuffix(filepath.Base(file), ".gitignore"))
		list[name] = data
	}

	return
}

func walker(path string, f os.FileInfo, err error) error {
	if filepath.Ext(path) == ".gitignore" {
		fileList = append(fileList, path)
	}

	return nil
}