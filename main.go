package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.homelessshelterdirectory.org"))

	c.OnHTML("div.layout_post_2", func(h *colly.HTMLElement) {
		// loops through the child elements of the matched HTML element 
		h.ForEach("img, h3, p, a", func(_ int, el *colly.HTMLElement) {
			fmt.Println(el.Text)
		})
	})

    c.Visit("https://www.homelessshelterdirectory.org/city/ky-lexington")
}
