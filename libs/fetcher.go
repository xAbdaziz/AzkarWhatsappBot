package libs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Zikr struct {
	Content string `json:"content"`
}

type PrayerTimes struct {
	Code int `json:"code"`
	Data struct {
		Timings struct {
			Fajr      time.Time `json:"Fajr"`
			Dhuhr     time.Time `json:"Dhuhr"`
			Asr       time.Time `json:"Asr"`
			Maghrib   time.Time `json:"Maghrib"`
			Isha      time.Time `json:"Isha"`
			LastThird time.Time `json:"Lastthird"`
		} `json:"timings"`
		Meta struct {
			Timezone string `json:"timezone"`
		} `json:"meta"`
	} `json:"data"`
}

func Fetch(URL string) ([]byte, error) {

	// Make an HTTP GET request to the URL
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Read the response body and return it
	return io.ReadAll(resp.Body)
}

func FetchZikr(ZikrCode string) string {
	// Check https://github.com/nawafalqari/azkar-api to see Zikr codes

	body, err := Fetch("https://azkar-api.nawafhq.repl.co/zekr?" + ZikrCode)
	if err != nil {
		panic(err)
	}

	var zikr Zikr
	err = json.Unmarshal(body, &zikr)
	if err != nil {
		panic(err)
	}

	return zikr.Content
}

func FetchPrayerTimes() PrayerTimes {
	city := os.Getenv("CITY")
	country := os.Getenv("COUNTRY")
	body, _ := Fetch(fmt.Sprintf("https://api.aladhan.com/v1/timingsByCity?city=%s&country=%s&iso8601=true&midnightMode=1", city, country))

	var pt PrayerTimes
	err := json.Unmarshal(body, &pt)
	if err != nil {
		panic(err)
	}

	return pt
}
