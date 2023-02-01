package main

import (
	"github.com/gocolly/colly"
)

func main() {
    c := colly.NewCollector(colly.AllowedDomains("www.homelessshelterdirectory.org"))

    c.OnHTML("")


    c.Visit("https://www.homelessshelterdirectory.org/city/ky-lexington")
}