package libs

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	"os"
	"time"
)

func Start(client *whatsmeow.Client) {
	bot := BotClient(client)
	pt := FetchPrayerTimes()
	if pt.Code != 200 {
		fmt.Println("Couldn't fetch Prayer times, did you correctly fill config.env?")
		os.Exit(1)
		return
	}

	city := os.Getenv("CITY")
	loc, _ := time.LoadLocation(pt.Data.Meta.Timezone)

	// Create a ticker that will run the code every minute.
	// Taken from https://gist.github.com/josephspurrier/ec57821bc4a3442a74ca
	t := MinuteTicker()

	for {
		// Wait for ticker to send
		<-t.C

		// Update the ticker
		t = MinuteTicker()

		// Time in location
		locTime := time.Now().In(loc)

		// Update timings everyday
		if locTime.Hour() == 0 && locTime.Minute() == 1 {
			pt = FetchPrayerTimes()
		}

		// Fajr
		if locTime.Hour() == pt.Data.Timings.Fajr.Hour() && locTime.Minute() == pt.Data.Timings.Fajr.Minute() {
			bot.sendMessage(fmt.Sprintf("حان وقت صلاة الفجر حسب توقيت مدينة %s", city))
		}

		// Dhuhur
		if locTime.Hour() == pt.Data.Timings.Dhuhr.Hour() && locTime.Minute() == pt.Data.Timings.Dhuhr.Minute() {
			prayerName := "الظهر"
			if pt.Data.Timings.Dhuhr.Weekday() == 5 {
				prayerName = "الجمعة"
			}
			bot.sendMessage(fmt.Sprintf("حان وقت صلاة %s حسب توقيت مدينة %s", prayerName, city))
		}

		// Asr
		if locTime.Hour() == pt.Data.Timings.Asr.Hour() && locTime.Minute() == pt.Data.Timings.Asr.Minute() {
			bot.sendMessage(fmt.Sprintf("حان وقت صلاة العصر حسب توقيت مدينة %s", city))
		}

		// Maghrib
		if locTime.Hour() == pt.Data.Timings.Maghrib.Hour() && locTime.Minute() == pt.Data.Timings.Maghrib.Minute() {
			bot.sendMessage(fmt.Sprintf("حان وقت صلاة المغرب حسب توقيت مدينة %s", city))
		}

		// Isha
		if locTime.Hour() == pt.Data.Timings.Isha.Hour() && locTime.Minute() == pt.Data.Timings.Isha.Minute() {
			bot.sendMessage(fmt.Sprintf("حان وقت صلاة العشاء حسب توقيت مدينة %s", city))
		}

		// Witr
		if locTime.Hour() == pt.Data.Timings.LastThird.Hour() && locTime.Minute() == pt.Data.Timings.LastThird.Minute() {
			bot.sendMessage(fmt.Sprintf("الوتر أحبتي ❤️"))
		}

		// Send Zikr every 2 hours
		if (locTime.Hour()%4 == 0 || locTime.Hour()%4 == 2) && locTime.Minute() == 0 {
			bot.sendMessage(FetchZikr("t"))
		}
	}
}
