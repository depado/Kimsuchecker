package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Depado/notificator"
)

var notify *notificator.Notificator
var availabilityUrl = "https://ws.ovh.com/dedicated/r2/ws.dispatcher/getAvailability2"

type OvhResponse struct {
	Answer struct {
		Class        string `json:"__class"`
		Availability []struct {
			Class            string `json:"__class"`
			DisplayMetazones int    `json:"displayMetazones"`
			Reference        string `json:"reference"`
			MetaZones        []struct {
				Class        string `json:"__class"`
				Availability string `json:"availability"`
				Zone         string `json:"zone"`
			} `json:"metaZones"`
			Zones []struct {
				Class        string `json:"__class"`
				Availability string `json:"availability"`
				Zone         string `json:"zone"`
			} `json:"zones"`
		} `json:"availability"`
	} `json:"answer"`
	Version string      `json:"version"`
	Error   interface{} `json:"error"`
	ID      int         `json:"id"`
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func notifyAvailability() {
	notify := notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "KS-3 Available",
	})

	err := notify.Push("KS-3 Available", "._.", "")
	perror(err)
}

func checkAvailability() {
	res, err := http.Get(availabilityUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	var data OvhResponse
	dec.Decode(&data)
	flag := false
	for _, res := range data.Answer.Availability {
		if res.Reference == "150sk30" {
			for _, zone := range res.Zones {
				if zone.Availability != "unknown" {
					flag = true
				}
			}
		}
	}
	if flag {
		notifyAvailability()
	}
}

func main() {
	for {
		go checkAvailability()
		time.Sleep(1 * time.Minute)
	}
}
