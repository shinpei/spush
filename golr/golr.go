package golr

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SolrConnector struct {
	host string
	port int
}

type SolrAddOption struct {
	Concurrency     int
	RecieverChannel chan []byte
}

// Assumes it'll get arrays of some data structure
func (sc *SolrConnector) AddDocuments(container interface{}, opt *SolrAddOption) {

	var err error
	// todo: size constrain should be placed here
	defer func() {
		if err != nil {
			log.Printf("Error occured, uploading document failed")
		}
	}()

	b, err := json.Marshal(container)
	if err != nil {
		log.Println("Failed at marshaling json structure, ", err)
	}
	respB, err := PostUpdate(sc.host,
		sc.port,
		b)
	if err != nil {
		log.Println(err)
	}
	opt.RecieverChannel <- respB
}

func Connect(host string, port int) *SolrConnector {
	return &SolrConnector{host, port}
}

func PostUpdate(host string, port int, payload []byte) ([]byte, error) {

	client := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/solr/update/json?commit=true", host, port)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Add("Content-type", "application/json")

	//	dump, _ := httputil.DumpRequestOut(req, true)
	//	fmt.Printf("%s", dump)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Recieved %d bytes.\n", len(body))
	return body, nil
}

type XMLNodeWalker interface {
	Walk(sc *SolrConnector, opt *SolrAddOption, decoder *xml.Decoder)
}

func (sc *SolrConnector) UploadXMLFile(
	path string,
	walker XMLNodeWalker,
	opt *SolrAddOption) {
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	walker.Walk(sc, opt, decoder)

}
