package main

import (
	"flag"
	"github.com/shinpei/spush/golr"
	"io/ioutil"
	"net/http"
	"runtime"
)

var inputFile = flag.String("infile", "jawiki-latest-pages-articles.xml", "Input file path")

type Data struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	TextCount int    `json:"text_count"`
}

func main() {
	con := golr.Connect("localhost", 8983)

	d := []Data{{
		Id:        "hige3",
		Title:     "hoge",
		Text:      "fuga",
		TextCount: 12,
	},
	}
	opt := &golr.SolrAddOption{
		Concurrency: runtime.NumCPU(),
	}
	con.AddDocuments(d, opt)

	//con.AddJSONFile(myjson, opt)
	con.AddXMLFile(*inputFile, opt)
}

func Get(url string) ([]byte, error) {
	r, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
