package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gocolly/colly"
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


var shelters []Shelter
func sheltersPage(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	city := r.FormValue("city")
	url := fmt.Sprintf("https://www.homelessshelterdirectory.org/city/%s-%s", state, city)
	fmt.Println("State:", state)
	fmt.Println("City:", city)

	// Create a new buffered writer
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	// Create a new Collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.homelessshelterdirectory.org"),
	)

	// OnHTML callback for each shelter
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
		if err := tmpl.Execute(buf, shelters); err != nil {
			log.Fatal(err)
		}
	})

	// Start the scraping process
	if err := c.Visit(url); err != nil {
		log.Printf("Error visiting %s: %s", url, err.Error())
	}
}

var foodPantries []FoodPantry
func foodPage(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	city := r.FormValue("city")
	url := fmt.Sprintf("https://www.foodpantries.org/ci/%s-%s", state, city)
	fmt.Println("State:", state)
	fmt.Println("City:", city)

	// Create a new buffered writer
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	// Create a new Collector
	c := colly.NewCollector(
		colly.AllowedDomains("https://www.foodpantries.org"),
	)

	// OnHTML callback for each shelter
	c.OnHTML("div.layout_post_2", func(e *colly.HTMLElement) {
		// Create a new Shelter instance and set its fields
		foodPantry := FoodPantry{
			Image:       e.ChildAttr("img", "src"),
			Name:        e.ChildText("h2 a"),
			Description: e.ChildText("div p"),
			SeeMore:     e.ChildAttr("a[href*=li]", "href"),
		}		

		// Append the pantry to the list
		foodPantries = append(foodPantries, foodPantry)
	})

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
		if err := tmpl.Execute(buf, foodPantries); err != nil {
			log.Fatal(err)
		}
	})

	// Start the scraping process
	if err := c.Visit(url); err != nil {
		log.Printf("Error visiting %s: %s", url, err.Error())
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

	http.HandleFunc("/information.html", sheltersPage)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

