package lyric

import (
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// GetLyric should receive url that was returned from GetLyricOptions. GetLyric
// returns lyric of the queried song.
func GetLyric(url string) (string, error) {

	var lyric string
	c := colly.NewCollector()

	c.OnHTML("span#ctl00_ContentPlaceHolder1_lblSubtitle", func(e *colly.HTMLElement) {
		content, err := e.DOM.Html()
		if err != nil {
			panic(err)
		}

		lyric = cleanHTML(content)
	})

	err := c.Visit(url + "&type=srt")
	if err != nil {
		return "", err
	}

	return lyric, nil
}

// GetLyricOptions queries available song lyrics. It returns map of title and
// url of the lyric.
func GetLyricOptions(search string) (map[string]string, error) {

	result := make(map[string]string)
	c := colly.NewCollector()

	c.OnHTML("#tablecontainer td a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		title := strings.TrimSpace(e.Text)
		result[title] = link
	})

	query := url.QueryEscape(search)
	err := c.Visit("https://www.rentanadviser.com/en/subtitles/subtitles4songs.aspx?src=" + query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// cleanHTML parses html text to valid utf-8 text
func cleanHTML(input string) string {

	content := html.UnescapeString(input)
	content = strings.ReplaceAll(content, "<br/>", "\n")
	// delete heading tag
	re := regexp.MustCompile(`^<h3>.*`)
	content = re.ReplaceAllString(content, "")
	// remove non-utf8 character
	re = regexp.MustCompile(`‚`)
	content = re.ReplaceAllString(content, ",")
	content = strings.ToValidUTF8(content, " ")
	content = strings.Map(func(r rune) rune {
		if r == 160 {
			return 32
		}
		return r
	}, content)

	return content
}