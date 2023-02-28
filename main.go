package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly"
)

type Shelter struct {
	Image      string
	Name       string
	Description string
	SeeMore string
}

func informationPage () {
	c := colly.NewCollector(colly.AllowedDomains("www.homelessshelterdirectory.org"))
	state := os.Args[1]
    city := os.Args[2]
    url := fmt.Sprintf("https://www.homelessshelterdirectory.org/city/%s-%s", state, city)
	var shelters []Shelter

	c.OnHTML("div.layout_post_2", func(h *colly.HTMLElement) {
		// create a new Shelter struct and set its fields based on the scraped data
		shelter := Shelter{
			Image: h.ChildAttr("img", "src"),
			Description:  h.ChildText("p"),
			Name:    h.ChildText("h4"),
			SeeMore: h.ChildAttr("a.btn_red", "href"),
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
	

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	  if r.Method == "POST" {
		http.Redirect(w, r, "/resources.html", http.StatusSeeOther)
		return
	  }
	  http.ServeFile(w, r, "./templates/home.html")
	})
  
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
  
	http.HandleFunc("/resources.html", func(w http.ResponseWriter, r *http.Request) {
	  http.ServeFile(w, r, "./templates/resources.html")
	})
  
	log.Fatal(http.ListenAndServe(":8080", nil))
  }
  

