package main

import (
	"flag"
	"fmt"
	"github.com/shinpei/spush/golr"
	"runtime"
)

var inputFilePath = flag.String("f", "jawiki-latest-pages-articles.xml", "Input file path")
var hostnameFlag = flag.String("h", "localhost", "Input solr server")
var portFlag = flag.Int("p", 8983, "Input port number")

func main() {

	flag.Parse()

	title := "spush"
	textBody := "spush is a tool for pushing documents to solr"
	con := golr.Connect(*hostnameFlag, *portFlag)
	d := []Page{{
		Id:        "spush",
		Title:     title,
		Text:      textBody,
		TextCount: len(textBody),
	},
	}

	opt := &golr.SolrAddOption{
		Concurrency: runtime.NumCPU(),
	}
	msg := <-con.AddDocuments(d, opt)
	fmt.Println(string(msg[:]))

	wikiWalker := &WikipediaXMLWalker{}
	con.UploadXMLFile(*inputFilePath, wikiWalker, opt)
}
