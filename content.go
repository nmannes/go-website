package main

import (
	"io"
	"math/rand"
	"text/template"

	"github.com/labstack/echo"
)

type Page struct {
	ImgInfo
	PageContent
}

type PageContent struct {
	Subhead        string
	ShowSubcontent bool
	Subcontent     string
	Content        []string
}

var masterTemplate *template.Template
var Images = []ImgInfo{}
var pages = map[string]PageContent{
	"/about": {
		Subhead:        "About",
		ShowSubcontent: true,
		Subcontent:     `Hi, my name is Nathan Mannes. I write <a href="https://golang.org">Go</a> at <a href="https://sezzle.com">Sezzle</a>. I grew up in New York City. I live in Minneapolis.`,
		Content: []string{
			`My tech work experience is detailed on my <a href="/resume">resume</a>`,
			`I graduated from Carleton College in 2019 with a degree in Computer Science`,
			`I like gaming (smash, mostly), tennis, musical theater, writing, non-fiction, and programming computers.`,
		},
	},
	"/news": {
		Subhead: "Nathan in the news",
		Content: []string{
			`In 2012, I was quoted in a local news article about NYC getting back to normal after <a href="https://www.cbsnews.com/news/nyc-area-schools-return-to-life-post-sandy/">
			hurricane Sandy</a>`,
			` In 2008, a piece of music that I wrote was <a href="https://www.nytimes.com/2008/01/14/arts/music/14youn.html">
			reviewed</a>`,
		},
	},
	"/history": {
		Subhead: "Major (and Minor) life events",
		Content: []string{

			`My grandma teaches me how to play tennis in 2003`,

			`In 2006 I am cast for a minor role in my <a href="https://www.schools.nyc.gov/schools/M199">elementary
			school</a>'s production of Pinnochio. I have no lines. It is a great success`,

			`In 2013, I first take a coding class at my <a href="https://stuy.enschool.org">high school</a>.
			It culminates in my creation of a reversi bot that is good enough to beat my dad`,

			`In the academic year of 2017-2018 I take 3 geology
			classes. I know more about rocks than I ever thought I wanted to`,

			`It is February of 2018. I get a call from an HR person at <a href="http://factset.com">Factset</a>.
			It turns out I did not blow the onsite interview. I accept this offer over the phone. I have gotten my first legit tech job`,

			` It is June 22nd, 2018. I am on a flight to Cleveland to visit my grandparents. I did not bring anything to do on the flight. I am sitting next to my brother. He is reading <a href="https://www.goodreads.com/book/show/1111.The_Power_Broker">this book</a>.
			He tells me to read the first chapter. I am happy. I have found something to do on the flight. I finish the book 6 months later`,

			`I graduate college June 8th, 2019`,

			`September 1st, 2019 I move to Minneapolis`,

			`September 11th, 2019 I start my second legit tech job at <a href="https://sezzle.com">Sezzle</a>`,

			`October 8th, 2019 I begin a tradition of bowling every Tuesday night with a few of my friends at <a href="https://www.bryantlakebowl.com">my local bowling alley</a>`,

			`Bowling night is on hiatus as of March 10th, 2020`,
		},
	},
	"/links": {

		Subhead: "links",
		Content: []string{
			`<a href="/resume">resume</a>`,
			`<a href="https://github.com/nmannes">github</a>`,
			`<a href="https://www.goodreads.com/user/show/48641482-nathan-mannes">goodreads</a>`,
			`<a href="https://linkedin.com/in/nathan-mannes">linkedin</a>`,
			`<a href="https://github.com/nmannes/go-website">the code for this website</a>`,
		},
	},
}

func Route(e echo.Context) error {

	path := e.Request().URL.Path

	pageContent, ok := pages[path]
	if !ok {
		pageContent = pages["/about"]
	}

	return RenderPage(e.Response().Writer, masterTemplate, Page{
		PageContent: pageContent,
		ImgInfo:     Images[rand.Intn(len(Images))],
	})
}

func RenderPage(w io.Writer, t *template.Template, p Page) error {

	err := t.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}
