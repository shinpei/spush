package golr

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"os"
	"regexp"
	"strconv"
	"strings"
)

var filter, _ = regexp.Compile("^%E3%83%95%E3%82%A1%E3%82%A4%E3%83%AB%3A.*|^help%3A.*|^talk%3A.*|^special%3A.*|^wikipedia%3A.*|^wikionary%3A.*|^user%3A.*|^user_talk%3A.*|^portal%3A.*|^mediawiki%3A.*|^template%3A.*|^Category%3A.*")

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
	fmt.Printf("Received %d bytes, ", len(body))
	return body, nil
}

//TODO: App-specific structure will be removed

type Page struct {
	Id        string `xml:"id" json:"id"`
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
	var pa []Page = make([]Page, opt.Concurrency*500)
	idx := 0
	var total int64 = 0
	var pushed int64 = 0
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
				p.Title = CannoTitle(p.Title)
				m := filter.MatchString(p.Title)
				if !m {
					p.Title, _ = url.QueryUnescape(p.Title)
					p.TextCount = len(p.Text)
					total++
					pa[idx] = p
					idx++
					fmt.Println("Added " + strconv.FormatInt(total, 10) + "/" + strconv.FormatInt(pushed, 10) + " for now..")
				} else {
					//println(p.Title)
				}
				pushed++
			}
			break
		default:
			break
		}
		if idx == opt.Concurrency*500-1 {
			println("hi")
			sc.AddDocuments(pa, opt)
			println("bye")
			idx = 0
		}
	}

}

func CannoTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, " ", "_", -1)
	can = url.QueryEscape(can)
	return can
}
