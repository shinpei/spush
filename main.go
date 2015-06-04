package main

import (
	"github.com/shinpei/spush/golr"
	"io/ioutil"
	"net/http"
	"runtime"
)

type Data struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	TextCount int    `json:"text_count"`
}

func main() {
	con := golr.Connect("localhost", 8983)

	d := []Data{{
		Id:        "hige2",
		Title:     "hoge",
		Text:      "fuga",
		TextCount: 12,
	},
	}
	opt := &golr.SolrAddOption{
		Concurrency: runtime.NumCPU(),
	}
	con.AddDocument(d, opt)

	//con.AddJSONFile(myjson, opt)
	//con.AddXMLFile(myxml, opt)

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
