package services

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type IManga struct {
	Name string
	URL  string
}

type INewMangaEntry struct {
	Thumbnail   string
	Name        string
	Description string
	Chapters    []IManga
}

var c = colly.NewCollector()

func GetMangasList() []IManga {
	var mangas []IManga

	c.OnHTML("div.post-title", func(e *colly.HTMLElement) {
		manga := IManga{
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

func NewMangaEntry(url string) INewMangaEntry {
	var thumbnail string
	var name string
	var description string
	var chapters []IManga

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

		chapters = append(chapters, IManga{
			Name: strings.Replace(element.Find("a").Text(), "\n", "", -1),
			URL:  link,
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

	return INewMangaEntry{
		Name:        name,
		Description: description,
		Thumbnail:   thumbnail,
		Chapters:    chapters,
	}
}

func GetMangaChapters(url string) []IManga {
	var chapters []IManga

	c.OnHTML("li.wp-manga-chapter", func(e *colly.HTMLElement) {
		element := e.DOM
		link, _ := element.Find("a").Attr("href")

		chapters = append(chapters, IManga{
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
