package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)
	
type WebhookBody struct {
	Embeds []EmbedBody `json:"embeds"`
}

type EmbedBody struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Image Image `json:"image"`
	Footer Footer `json:"footer"`
}

type Image struct {
	URL string `json:"url"`
}

type Footer struct {
	Text string `json:"text"`
}

func sendWebhook(WebhookURL string, Basin NHCBasinItem, datetime time.Time) {
	data := &WebhookBody{
		Embeds: []EmbedBody{
			{
				Title: "NHC Atlantic Outlook",
				URL: Basin.Link,
				Image: Image{
					URL: "https://www.nhc.noaa.gov/archive/xgtwo/atl/" + formatDate(datetime) + "/two_atl_5d0.png",
				},
				Footer: Footer{ 
					Text: "Last updated " + datetime.Local().Format("Jan 02 2006, 3:04pm MST"),
				},
			},
		},
	}

	dataStr, _ := json.Marshal(data)

	resp, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(dataStr))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if resp.StatusCode != 204 {
		fmt.Fprintln(os.Stderr, "Error while sending Webhook message,", resp.Status, string(body))
		return
	}

	fmt.Println("Webhook message has been sent!")
}

func formatDate(datetime time.Time) string {
	return strings.Join(strings.Split(datetime.Format("2006 01 02 15 04"), " "), "")
}