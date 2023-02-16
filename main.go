package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.homelessshelterdirectory.org"))

	c.OnHTML("div.layout_post_2.clearfix", func(h *colly.HTMLElement) {
		fmt.Println(h.Text)
	})

	c.Visit("https://www.homelessshelterdirectory.org/city/ky-lexington")
}