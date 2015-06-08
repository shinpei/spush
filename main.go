package main

import (
	"flag"
	"github.com/shinpei/spush/golr"
	"io/ioutil"
	"net/http"
	"runtime"
)

var inputFilePath = flag.String("infile", "jawiki-latest-pages-articles.xml", "Input file path")

type Data struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	TextCount int    `json:"text_count"`
}

//TODO: App-specific structure will be removed

func main() {
	con := golr.Connect("localhost", 8983)

	d := []Data{{
		Id:        "hige3",
		Title:     "hoge",
		Text:      "fuga",
		TextCount: 12,
	},
	}

	recvChan := make(chan []byte)
	opt := &golr.SolrAddOption{
		Concurrency:     runtime.NumCPU(),
		ReceiverChannel: recvChan,
	}
	go con.AddDocuments(d, opt)
	msg := <-recvChan
	println(string(msg[:]))

	wikiWalker := &WikipediaXMLWalker{}
	con.UploadXMLFile(*inputFilePath, wikiWalker, opt)

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
