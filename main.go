package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("NOAA NHC Atlantic Basin Tracker")

	if !checkEnv() {
		fmt.Fprintln(os.Stderr, "Missing environment variables.")
		return
	}

	lastPublish, err := getLastPublishDate()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	resp, err := http.Get("https://www.nhc.noaa.gov/gtwo.xml")

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	var data NHCData

	xml.Unmarshal(body, &data)

	items := data.Channel.Item

	var basin NHCBasinItem

	for i := 0; i < len(items); i++ {
		if items[i].Title == "NHC Atlantic Outlook" {
			basin = items[i]
			break
		}
	}

	if basin.Title == "" {
		fmt.Println("Error finding Atlantic Basin Tropical Outlook.")
		return
	}

	layout := "Mon, 02 Jan 2006 15:04:05 MST"
	datetime, err := time.Parse(layout, basin.PubDate)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Println("Atlantic Basin Tropical Outlook was last published on", datetime.Local().Format("Jan 02 2006, 3:04pm MST"))

	previewURL := "https://www.nhc.noaa.gov/archive/xgtwo/atl/" + formatDate(datetime) + "/two_atl_5d0.png"

	if lastPublish != datetime.Unix()*1000 {
		if !verifyPreviewAvailable(previewURL) {
			fmt.Println("New outlook information available but preview is not present.")
			return
		}

		err := setLastPublishDate(datetime.Unix() * 1000)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		sendWebhook(os.Getenv("WEBHOOK_URL"), basin, datetime, previewURL)
	} else {
		fmt.Println("No new outlook information present, check again next time!")
	}
}

func checkEnv() bool {
	return os.Getenv("WEBHOOK_URL") != "" &&
		os.Getenv("REDIS_DSN") != ""
}

func formatDate(datetime time.Time) string {
	return strings.Join(strings.Split(datetime.Format("2006 01 02 15 04"), " "), "")
}

func verifyPreviewAvailable(PreviewURL string) bool {
	resp, err := http.Get(PreviewURL)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	}

	return resp.StatusCode == 200
}
