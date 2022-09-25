package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"
)

type WebhookBody struct {
	Embeds []EmbedBody `json:"embeds"`
}

type EmbedBody struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Image  Image  `json:"image"`
	Footer Footer `json:"footer"`
}

type Image struct {
	URL string `json:"url"`
}

type Footer struct {
	Text string `json:"text"`
}

func sendWebhook(WebhookURL string, Item NHCItem, datetime time.Time, PreviewURL string, walletId string) {
	data := &WebhookBody{
		Embeds: []EmbedBody{
			{
				Title: Item.Title,
				URL:   Item.Link,
				Image: Image{
					URL: PreviewURL,
				},
				Footer: Footer{
					Text: "Last updated " + datetime.Local().Format("Jan 02 2006, 3:04pm MST"),
				},
			},
		},
	}

	dataStr, _ := json.Marshal(data)

	resp, err := client.Post(WebhookURL, "application/json", bytes.NewBuffer(dataStr))

	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	if resp.StatusCode != 204 {
		log.Println("Error while sending Webhook message,", resp.Status, string(body))
		return
	}

	log.Println("Webhook message has been sent!", walletId)
}
