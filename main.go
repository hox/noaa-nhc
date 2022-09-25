package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var client http.Client

func main() {
	client = http.Client{
		Timeout: 10 * time.Second,
	}

	log.Println("NOAA NHC Atlantic Basin Tracker")

	if !checkEnv() {
		log.Println("Missing environment variables.")
		return
	}

	fetchTropicalOutlook()

	if os.Getenv("WALLET_ID") != "" {
		log.Println()
		fetchWallet(os.Getenv("WALLET_ID"))
	}
}

func fetchWallet(walletId string) {
	log.Printf("Fetching Atlantic Basin Tropical Wallet #%s\n", walletId)

	data, err := fetchNHC(fmt.Sprintf("https://www.nhc.noaa.gov/nhc_at%s.xml", walletId))

	if err != nil {
		log.Println(err.Error())
		return
	}

	advisories := data.Channel.Item

	var summary NHCItem

	for i := range advisories {
		if strings.HasPrefix(advisories[i].Title, "Summary") {
			summary = advisories[i]
			break
		}
	}

	validateAndSendPublication(summary, &walletId)
}

func fetchTropicalOutlook() {
	data, err := fetchNHC("https://www.nhc.noaa.gov/gtwo.xml")

	if err != nil {
		log.Println(err.Error())
		return
	}

	basins := data.Channel.Item

	var basin NHCItem

	for i := range basins {
		if basins[i].Title == "NHC Atlantic Outlook" {
			basin = basins[i]
			break
		}
	}

	if basin.Title == "" {
		log.Println("Error finding Atlantic Basin Tropical Outlook.")
		return
	}

	validateAndSendPublication(basin, nil)
}

func validateAndSendPublication(item NHCItem, walletId *string) {
	lastPublish, err := getLastPublishDate(walletId)

	if err != nil {
		log.Println(err.Error())
		return
	}

	layout := "Mon, 02 Jan 2006 15:04:05 MST"
	datetime, err := time.Parse(layout, item.PubDate)

	if err != nil {
		log.Println(err.Error())
		return
	}

	var previewURL string

	id := ""

	if walletId == nil {
		log.Println("Atlantic Basin Tropical Outlook was last published on", datetime.Local().Format("Jan 02 2006, 3:04pm MST"))
		previewURL = fmt.Sprintf("https://www.nhc.noaa.gov/archive/xgtwo/atl/%s/two_atl_5d0.png", formatDate(datetime, false))
	} else {
		id = "#" + *walletId
		log.Printf("Atlantic Basin Tropical Wallet #%s was last published on %s\n", *walletId, datetime.Local().Format("Jan 02 2006, 3:04pm MST"))
		alId := strings.ToUpper(strings.Split(item.Guid, "-")[1])
		shortId := strings.Join(strings.Split(alId, "")[2:4], "")
		datestr := formatDate(datetime, true)
		previewURL = fmt.Sprintf("https://www.nhc.noaa.gov/storm_graphics/AT%s/refresh/%s_5day_cone_no_line_and_wind+png/%s.png", shortId, alId, datestr)
	}

	if lastPublish != datetime.Unix()*1000 {
		if !verifyPreviewAvailable(previewURL) {
			log.Println("New outlook information available but preview is not present.", id)
			return
		}

		if err := setLastPublishDate(datetime.Unix()*1000, walletId); err != nil {
			log.Println(err.Error())
			return
		}

		sendWebhook(os.Getenv("WEBHOOK_URL"), item, datetime, previewURL, id)
	} else {
		log.Println("No new outlook information present, check again next time!", id)
	}
}
