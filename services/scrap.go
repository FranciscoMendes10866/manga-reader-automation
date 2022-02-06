package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/types"
	"github.com/gocolly/colly"
)

var c = colly.NewCollector(colly.MaxDepth(1), colly.DetectCharset(), colly.AllowURLRevisit())

func GetMangasList() []types.IManga {
	var mangas []types.IManga

	c.SetProxyFunc(config.SetScrappingProxy())
	c.SetRequestTimeout(120 * time.Second)
	c.Limit(helpers.ScrapLimitOptions)

	c.OnHTML("div.post-title", func(e *colly.HTMLElement) {
		manga := types.IManga{
			Name: strings.Replace(e.Text, "\n", "", -1),
			URL:  e.ChildAttr("a[href]", "href"),
		}
		mangas = append(mangas, manga)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://mangaclash.com/")

	return mangas
}

func NewMangaEntry(url string) types.INewMangaEntry {
	var thumbnail string
	var name string
	var description string
	var chapters []types.IManga
	var categories []string

	c.SetProxyFunc(config.SetScrappingProxy())
	c.SetRequestTimeout(120 * time.Second)
	c.Limit(helpers.ScrapLimitOptions)

	c.OnHTML("div.post-title", func(e *colly.HTMLElement) {
		name = strings.Replace(e.Text, "\n", "", -1)
	})

	c.OnHTML("div.summary__content", func(e *colly.HTMLElement) {
		description = strings.Replace(e.Text, "\n", "", -1)
	})

	c.OnHTML("div.genres-content", func(e *colly.HTMLElement) {
		element := e.DOM
		category := element.Text()
		var splited = strings.Split(category, ",")

		for _, value := range splited {
			categories = append(categories, strings.TrimSpace(value))
		}
	})

	c.OnHTML("div.summary_image", func(e *colly.HTMLElement) {
		element := e.DOM
		link, _ := element.Find("img").Attr("data-src")
		thumbnail = link
	})

	c.OnHTML("li.wp-manga-chapter", func(e *colly.HTMLElement) {
		element := e.DOM
		link, _ := element.Find("a").Attr("href")

		chapters = append(chapters, types.IManga{
			Name: strings.Replace(element.Find("a").Text(), "\n", "", -1),
			URL:  link,
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

	return types.INewMangaEntry{
		Name:        name,
		Description: description,
		Thumbnail:   thumbnail,
		Chapters:    chapters,
		Categories:  categories,
	}
}

func GetMangaChapters(url string) []types.IManga {
	var chapters []types.IManga

	c.SetProxyFunc(config.SetScrappingProxy())
	c.SetRequestTimeout(120 * time.Second)
	c.Limit(helpers.ScrapLimitOptions)

	c.OnHTML("li.wp-manga-chapter", func(e *colly.HTMLElement) {
		element := e.DOM
		link, _ := element.Find("a").Attr("href")

		chapters = append(chapters, types.IManga{
			Name: strings.Replace(element.Find("a").Text(), "\n", "", -1),
			URL:  link,
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

	return chapters
}

func GetChapterPages(url string) []string {
	var pages []string

	c.SetProxyFunc(config.SetScrappingProxy())
	c.SetRequestTimeout(120 * time.Second)
	c.Limit(helpers.ScrapLimitOptions)

	c.OnHTML("div.page-break", func(e *colly.HTMLElement) {
		element := e.DOM
		link, _ := element.Find("img").Attr("data-src")

		pages = append(pages, strings.TrimSpace(link))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

	return pages
}
