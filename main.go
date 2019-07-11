package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type rss2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	// Optional
	PubDate  string `xml:"channel>pubDate"`
	ItemList []item `xml:"channel>item"`
}

type item struct {
	// Required
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"`
	// Optional
	Content   template.HTML `xml:"encoded"`
	PubDate   string        `xml:"pubDate"`
	Comments  string        `xml:"comments"`
	Enclosure enclosure     `xml:"enclosure"`
}
type enclosure struct {
	Url string `xml:"url,attr"`
}

func main() {
	r := rss2{}
	res, err := http.Get("https://feed.megafono.host/benzina-no-meiao")
	if err != nil {
		log.Fatal("cant get feed")
	}
	defer res.Body.Close()

	xmlContent, _ := ioutil.ReadAll(res.Body)
	err = xml.Unmarshal(xmlContent, &r)
	if err != nil {
		panic(err)
	}

	os.MkdirAll("episodios", 0777)
	for _, item := range r.ItemList {

		if strings.Contains(item.Title, "|") == false {
			continue
		}
		fmt.Println(item.Title)
		fmt.Println(item.Enclosure.Url)
		resEp, err := http.Get(item.Enclosure.Url)
		if err != nil {
			log.Fatal(fmt.Sprintf("cant get %s", item.Title))
		}
		defer resEp.Body.Close()

		file, err := os.Create(fmt.Sprintf("episodios/%s.mp3", item.Title))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		io.Copy(file, resEp.Body)

	}
}
