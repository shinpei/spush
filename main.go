package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Data struct {
	Id    string
	Title string
}

func main() {

	b, err := json.Marshal(Data{"1", "mybook"})
	if err != nil {
		panic(err)
	}
	respB, err := Post("http://localhost:8983/update?wt=json", &b)
	if err != nil {
		panic(err)
	}
	var datas []Data
	err = json.Unmarshal(respB, &datas)
	if err != nil {
		s := string(respB[:])
		println(s)
	} else {
		fmt.Printf("%%x\n", datas[0])
	}
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

func Post(url string, payload *[]byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(*payload))
	req.Header.Add("Content-type", "application/json")
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
