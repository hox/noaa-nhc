package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func checkEnv() bool {
	return os.Getenv("WEBHOOK_URL") != "" &&
		os.Getenv("REDIS_DSN") != ""
}

func formatDate(datetime time.Time, short bool) string {
	var datestr string
	if short {
		datestr = "15 04 05"
	} else {
		datestr = "2006 01 02 15 04"
	}
	return strings.Join(strings.Split(datetime.Format(datestr), " "), "")
}

func fetchNHC(url string) (NHCData, error) {
	resp, err := client.Get(url)

	if err != nil {
		return NHCData{}, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return NHCData{}, err
	}

	var data NHCData

	xml.Unmarshal(body, &data)

	return data, nil
}

func verifyPreviewAvailable(PreviewURL string) bool {
	resp, err := client.Get(PreviewURL)

	if err != nil {
		log.Println(err.Error())
		return false
	}

	return resp.StatusCode == 200
}
