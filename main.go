package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

	if lastPublish == 0 {
		fmt.Println("\"noaa-nhc:last-updated\" key not found in redis, setting to current publish date and sending webhook.")
		err := setLastPublishDate(datetime.Unix()*1000)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	} else if lastPublish == datetime.Unix()*1000 {
		fmt.Println("No new outlook information available.")
		return
	} else {
		err := setLastPublishDate(datetime.Unix()*1000)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}

	sendWebhook(os.Getenv("WEBHOOK_URL"), basin, datetime)
}

func checkEnv() bool {
	return os.Getenv("WEBHOOK_URL") != "" &&
				 os.Getenv("REDIS_DSN") != ""
}