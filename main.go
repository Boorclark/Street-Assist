package main

import (
	"bufio"         // Provides buffered I/O operations
	"fmt"           // Provides formatted I/O operations
	"html/template" // Provides tools for rendering HTML templates
	"log"           // Provides logging functions
	"net/http"      // Provides HTTP client and server implementations

	"github.com/gocolly/colly" // Scraping framework for Go
)

type Shelter struct {
	Image      string
	Name       string
	Description string
	SeeMore string
}

type FoodPantry struct {
	Image		string
	Name		string
	Description	string
	SeeMore		string
}

var (
	shelters     []Shelter // Slice of shelter stuct
	foodPantries []FoodPantry // Slice of FoodPantry struct
	state        string
	city         string
)

func resourcesPage(w http.ResponseWriter, r *http.Request, state string, city string) {
	// Creates a empty slice of shelter and food pantry
	shelters = make([]Shelter, 0)  
	foodPantries = make([]FoodPantry, 0)

	// Extracts HTTP request path (/information/shelters, /information/food)
	path := r.URL.Path 

	// Formats and returns URL strings
	shelterURL := fmt.Sprintf("https://www.homelessshelterdirectory.org/city/%s-%s", state, city) 
	foodURL := fmt.Sprintf("https://www.foodpantries.org/ci/%s-%s", state, city)

	// Create a new buffered writer. This reduces the number of system calls.
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	// Create a new Collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.homelessshelterdirectory.org", "www.foodpantries.org"),
	)

	// OnHTML callback for shelter information
	c.OnHTML("div.layout_post_2", func(e *colly.HTMLElement) {
		// Create a new Shelter instance and set its fields
		shelter := Shelter{
			Image:       e.ChildAttr("img", "src"),
			Name:        e.ChildText("h4"),
			Description: e.ChildText("p"),
			SeeMore:     e.ChildAttr("a.btn_red", "href"),
		}
		// Append the Shelter to the list
		shelters = append(shelters, shelter)
	})

	// OnHTML callback for food pantry information
	c.OnHTML(".blog-list h2", func(e *colly.HTMLElement) {
		// Create a new FoodPantry instance and set its fields
		foodPantry := FoodPantry{
			Image: e.DOM.Next().Next().AddBack().AttrOr("src", ""),
			Name:        e.ChildText("h2 a"),
			Description: e.DOM.Next().Next().Next().Text(),
			SeeMore:     e.ChildAttr("a", "href"),
		}
		// Append the FoodPantry to the list
		foodPantries = append(foodPantries, foodPantry)})
	// OnError callback to handle errors
	c.OnError(func(_ *colly.Response, err error) {
		log.Printf("Error scraping: %s", err.Error())
	})

	// OnScraped callback to execute once the scraping is done
	c.OnScraped(func(_ *colly.Response) {
		// Parse the information template
		tmpl, err := template.ParseFiles("templates/information.html")
		if err != nil {
			log.Fatal(err)
		}

		// Generate the HTML and write it to the buffered writer
		if path == "/information/shelters" {
			if err := tmpl.Execute(buf, shelters); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := tmpl.Execute(buf, foodPantries); err != nil {
				log.Fatal(err)
			}
		}
	})

	// Start the scraping process
	if path == "/information/shelters" {
		if err := c.Visit(shelterURL); err != nil {
			log.Printf("Error visiting %s: %s", shelterURL, err.Error())
		}
	} 
	if path == "/information/food" {
		if err := c.Visit(foodURL); err != nil {
			log.Printf("Error visiting %s: %s", foodURL, err.Error())
		}
	} 
}


func informationHandler(w http.ResponseWriter, r *http.Request) {
	// Check if state and city values have already been stored
	if state == "" || city == "" {
		state = r.FormValue("state")
		city = r.FormValue("city")
		fmt.Println("State:", state)
		fmt.Println("City:", city)
	}

	// Pass on the stored or extracted values to the resourcesPage function
	resourcesPage(w, r, state, city)
}

func main() {
	// Handler functions for URL's
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

	http.HandleFunc("/emergency-services", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/emergency.html")
	})

	http.HandleFunc("/information/", func(w http.ResponseWriter, r *http.Request) {
		informationHandler(w, r)
	})

	// Start the web server
	log.Fatal(http.ListenAndServe(":8080", nil))
}