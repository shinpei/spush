package main

import (
	"flag"
	"github.com/shinpei/spush/golr"
	"runtime"
)

var inputFilePath = flag.String("infile", "jawiki-latest-pages-articles.xml", "Input file path")
var hostnameFlag = flag.String("hostname", "localhost", "Input solr server")
var portFlag = flag.Int("port", 8983, "Input port number")

func main() {
	con := golr.Connect(*hostnameFlag, *portFlag)
	d := []Page{{
		Id:        "hige3",
		Title:     "hoge",
		Text:      "fuga",
		TextCount: 12,
	},
	}

	recvChan := make(chan []byte)
	opt := &golr.SolrAddOption{
		Concurrency:     runtime.NumCPU(),
		RecieverChannel: recvChan,
	}
	go con.AddDocuments(d, opt)
	msg := <-recvChan
	close(recvChan)

	println(string(msg[:]))

	wikiWalker := &WikipediaXMLWalker{}
	con.UploadXMLFile(*inputFilePath, wikiWalker, opt)

}
