package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type Shelter struct {
	ImageURL string
	Name     string
	Address  string
	Phone    string
}

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.homelessshelterdirectory.org"))

	var shelters []Shelter

	c.OnHTML("div.layout_post_2", func(h *colly.HTMLElement) {
		// create a new Shelter struct and set its fields based on the scraped data
		shelter := Shelter{
			ImageURL: h.ChildAttr("img", "src"),
			Name:     h.ChildText("h3"),
			Address:  h.ChildText("p"),
			Phone:    h.ChildText("a"),
		}
		// add the new shelter to the list of shelters
		shelters = append(shelters, shelter)
	})

	c.OnScraped(func(r *colly.Response) {
		tmpl, err := template.ParseFiles("templates/information.html")
		if err != nil {
			log.Fatal(err)
		}
	
		// Create a buffer to store the generated HTML
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, shelters)
		if err != nil {
			log.Fatal(err)
		}
	
		// Serve the generated HTML as the HTTP response
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write(buf.Bytes())
		})
	
		// Start the HTTP server and listen for incoming requests
		log.Fatal(http.ListenAndServe(":8080", nil))
	})
	

	err := c.Visit("https://www.homelessshelterdirectory.org/city/ky-lexington")
	if err != nil {
		log.Fatal(err)
	}
}
