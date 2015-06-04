package golr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type SolrConnector struct {
	host string
	port int
}

type SolrAddOption struct {
	Concurrency int
}

func (sc *SolrConnector) AddDocument(container interface{}, opt *SolrAddOption) {

	b, err := json.Marshal(container)
	ss := string(b[:])
	println(ss)
	if err != nil {
		panic(err)
	}

	respB, err := Post("http://localhost:8983/solr/update/json?commit=true", b)

	if err != nil {
		panic(err)
	}
	var datas interface{}
	err = json.Unmarshal(respB, &datas)

	if err != nil {
		s := string(respB[:])
		println(s)
	} else {
		fmt.Printf("%x\n", datas)
	}

}

func Connect(host string, port int) *SolrConnector {
	return &SolrConnector{host, port}
}
func Post(url string, payload []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Add("Content-type", "application/json")
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s", dump)
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
