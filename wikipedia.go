package main

import (
	"encoding/xml"
	"fmt"
	"github.com/shinpei/spush/golr"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var filter, _ = regexp.Compile("^%E3%83%95%E3%82%A1%E3%82%A4%E3%83%AB%3A.*|^help%3A.*|^talk%3A.*|^special%3A.*|^wikipedia%3A.*|^wikionary%3A.*|^user%3A.*|^user_talk%3A.*|^portal%3A.*|^mediawiki%3A.*|^template%3A.*|^category%3A.*|^wp%3A.*")

type Page struct {
	Id        string `xml:"id" json:"id"`
	Title     string `xml:"title" json:"title"`
	Text      string `xml:"revision>text" json:"text"`
	TextCount int    `json:"text_count"`
}

func CannoTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, " ", "_", -1)
	can = url.QueryEscape(can)
	return can
}

type WikipediaXMLWalker struct{}

func (w *WikipediaXMLWalker) Walk(inputChan chan interface{},
	opt *golr.SolrAddOption, decoder *xml.Decoder) {

	var inElement string
	PageChunk := 500
	var pa []Page = make([]Page, opt.Concurrency*PageChunk)
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
				}
				pushed++
			}
			break
		default:
			break
		}

		if idx == opt.Concurrency*PageChunk-1 {
			fmt.Println("Added " + strconv.FormatInt(total, 10) + "/" + strconv.FormatInt(pushed, 10) + " for now..")
			inputChan <- pa
			idx = 0
		}
	}
}
