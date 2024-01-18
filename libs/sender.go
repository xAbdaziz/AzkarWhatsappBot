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

	city := os.Getenv("CITY")
	loc, err := time.LoadLocation(pt.Data.Meta.Timezone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		locTime := time.Now().In(loc)
		if isNewDay(locTime) {
			pt = FetchPrayerTimes() // handle error appropriately
		}

		checkAndSendPrayerAlerts(*bot, pt, locTime, city)
		sendZikrEveryTwoHours(*bot, locTime)
	}
}

func isNewDay(t time.Time) bool {
	return t.Hour() == 0 && t.Minute() == 1
}

func checkAndSendPrayerAlerts(bot Bot, pt PrayerTimes, locTime time.Time, city string) {
	timings := []struct {
		Time             time.Time
		Name             string
		SpecialCondition func(time.Time) bool
	}{
		{pt.Data.Timings.Fajr, "الفجر", nil},
		{pt.Data.Timings.Dhuhr, "الظهر", func(t time.Time) bool { return t.Weekday() == time.Friday }},
		{pt.Data.Timings.Asr, "العصر", nil},
		{pt.Data.Timings.Maghrib, "المغرب", nil},
		{pt.Data.Timings.Isha, "العشاء", nil},
		{pt.Data.Timings.LastThird, "الوتر", nil},
	}

	for _, timing := range timings {
		if locTime.Hour() == timing.Time.Hour() && locTime.Minute() == timing.Time.Minute() {
			message := fmt.Sprintf("حان وقت صلاة %s حسب توقيت مدينة %s", timing.Name, city)
			if timing.SpecialCondition != nil && timing.SpecialCondition(locTime) {
				message = fmt.Sprintf("حان وقت صلاة الجمعة حسب توقيت مدينة %s", city)
			}
			bot.sendMessage(message)
		}
	}
}

func sendZikrEveryTwoHours(bot Bot, locTime time.Time) {
	if (locTime.Hour()%4 == 0 || locTime.Hour()%4 == 2) && locTime.Minute() == 0 {
		bot.sendMessage(FetchZikr())
	}
}
