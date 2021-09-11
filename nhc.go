package main

import "encoding/xml"

type NHCData struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Dc      string   `xml:"dc,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text string `xml:",chardata"`
		Link struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		PubDate        string `xml:"pubDate"`
		LastBuildDate  string `xml:"lastBuildDate"`
		Title          string `xml:"title"`
		Description    string `xml:"description"`
		Copyright      string `xml:"copyright"`
		ManagingEditor string `xml:"managingEditor"`
		Language       string `xml:"language"`
		WebMaster      string `xml:"webMaster"`
		Image          struct {
			Text        string `xml:",chardata"`
			URL         string `xml:"url"`
			Link        string `xml:"link"`
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Width       string `xml:"width"`
			Height      string `xml:"height"`
		} `xml:"image"`
		Item []NHCBasinItem `xml:"item"`
	} `xml:"channel"`
}

type NHCBasinItem struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	Author      string `xml:"author"`
}