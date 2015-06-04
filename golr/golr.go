package golr

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

var filter, _ = regexp.Compile("^file:.*|^talk:.*|^special:.*|^wikipedia:.*|^wikionary:.*|^user:.*|^user_talk:.*")

type SolrConnector struct {
	host string
	port int
}

type SolrAddOption struct {
	Concurrency int
}

// Assumes it'll get arrays of some data structure
func (sc *SolrConnector) AddDocuments(container interface{}, opt *SolrAddOption) {

	// todo: size constrain should be placed here
	b, err := json.Marshal(container)

	if err != nil {
		panic(err)
	}

	respB, err := PostUpdate(sc.host,
		sc.port,
		b)

	if err != nil {
		panic(err)
	}

	var datas interface{}
	err = json.Unmarshal(respB, &datas)

	if err != nil {
		s := string(respB[:])
		println(s)
		panic(err)
	} else {
		//		fmt.Printf("%x\n", datas)
	}

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
	fmt.Printf("Received %d bytes\n", len(body))
	return body, nil
}

//TODO: App-specific structure will be removed

type Page struct {
	Id        string `json:"id"`
	Title     string `xml:"title" json:"title"`
	Text      string `xml:"revision>text" json:"text"`
	TextCount int    `json:"text_count"`
}

func (sc *SolrConnector) AddXMLFile(path string, opt *SolrAddOption) {
	xmlFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var inElement string
	var pa []Page = make([]Page, opt.Concurrency)
	idx := 0
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "page" {
				var p Page
				decoder.DecodeElement(&p, &se)
				p.TextCount = len(p.Text)
				p.Id = string(time.Now().UnixNano())
				pa[idx] = p
				idx++
			}
			break
		default:
		}
		if idx == opt.Concurrency-1 {
			sc.AddDocuments(pa, opt)
			idx = 0
		}

	}
}
