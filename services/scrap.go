package services

import (
	"fmt"
	"strings"

	"github.com/FranciscoMendes10866/queues/types"
	"github.com/gocolly/colly"
)

var c = colly.NewCollector()

func GetMangasList() []types.IManga {
	var mangas []types.IManga

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

	c.OnHTML("div.post-title", func(e *colly.HTMLElement) {
		name = strings.Replace(e.Text, "\n", "", -1)
	})

	c.OnHTML("div.summary__content", func(e *colly.HTMLElement) {
		description = strings.Replace(e.Text, "\n", "", -1)
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
	}
}

func GetMangaChapters(url string) []types.IManga {
	var chapters []types.IManga

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
