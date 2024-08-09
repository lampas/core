package timeago

import (
	"time"

	tgo "github.com/SerhiiCho/timeago"
)

func Format(lang string, date time.Time) string {
	tgo.SetConfig(tgo.Config{
		Language: lang,
	})

	return tgo.Parse(date)
}
