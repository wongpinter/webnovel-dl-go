package webnovel

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Scraper struct {
	URL string
	Webnovel *Webnovel
}

func (s *Scraper) Fetch(uRL string) {
	s.URL = uRL

	collector := colly.NewCollector()
	chapterCollector := collector.Clone()

	contentCollector := colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.Async(true),
	)

	err := contentCollector.Limit(
		&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 2,
		})

	if err != nil {
		fmt.Println(err)
	}

	s.Webnovel = new(Webnovel)

	s.getWebnovelMetadata(collector, chapterCollector)

	s.getChapters(chapterCollector, contentCollector)

	_ = collector.Visit(s.URL)

	contentCollector.Wait()

	sort.Sort(s.Webnovel.Chapters)

	s.Webnovel.TotalChapter = len(s.Webnovel.List)

	s.Webnovel.Save()
}

func (s Scraper) getChapters(chapterCollector *colly.Collector, contentCollector *colly.Collector) {
	chapterCollector.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		chapterUrls := fmt.Sprintf("%s%s", `https://readnovelfull.com`, e.Attr(`href`))

		err := contentCollector.Visit(chapterUrls)
		if err != nil {
			fmt.Println(err)
		}
	})

	contentCollector.OnHTML(`div#chapter`, func(e *colly.HTMLElement) {
		chapter := &Chapter{}

		chapter.ID, _ = strconv.Atoi(e.ChildAttr(`button.btn-warning[data-chr-id]`, `data-chr-id`))
		chapter.Title = getTitle(e)
		chapter.Content = getContent(e)

		s.Webnovel.AddChapter(chapter)
	})

	contentCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Parsing .. ", r.URL.String())
	})
}

func (s *Scraper) getWebnovelMetadata(collector *colly.Collector, chapterCollector *colly.Collector) {
	collector.OnHTML(`div.col-xs-12.col-info-desc`, func(e *colly.HTMLElement) {
		s.Webnovel.Name = e.ChildText(`h3.title:last-child`)
		s.Webnovel.Author = e.ChildTexts(`ul.info.info-meta li:nth-child(2) a`)
		s.Webnovel.Genre = e.ChildTexts(`ul.info.info-meta li:nth-child(3) a`)
		s.Webnovel.URL = e.Request.URL.String()
	})

	collector.OnHTML(`div#rating`, func(e *colly.HTMLElement) {
		chapterURL := fmt.Sprintf("%s%s", `https://readnovelfull.com/ajax/chapter-archive?novelId=`, e.Attr(`data-novel-id`))

		err := chapterCollector.Visit(chapterURL)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func getTitle(element *colly.HTMLElement) string {
	return element.ChildText(`span.chr-text`)
}

func getContent(element *colly.HTMLElement) string {
	var content string

	element.ForEach(`p`, func(i int, e *colly.HTMLElement) {
		re := regexp.MustCompile(`^Chapter.+`)
		s := re.ReplaceAllString(e.Text, ``)

		s = strings.Trim(s, " ")

		content += fmt.Sprintf(`%s%s`, s, "\n\n")
	})

	return content
}